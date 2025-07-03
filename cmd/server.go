package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "github.com/JRaver/k8s-controller-tutorial/docs"
	"github.com/JRaver/k8s-controller-tutorial/pkg/api"
	frontendv1alpha1 "github.com/JRaver/k8s-controller-tutorial/pkg/apis/frontend/v1alpha1"
	"github.com/JRaver/k8s-controller-tutorial/pkg/ctrl"
	"github.com/JRaver/k8s-controller-tutorial/pkg/informer"
	"github.com/JRaver/k8s-controller-tutorial/pkg/telemetry"
	"github.com/buaazp/fasthttprouter"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrlruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

// @title K8s Controller Tutorial API
// @version 1.0
// @description My awesome lab controller with Swagger UI
// @host localhost:8080
// @BasePath /

var serverPort int = 8080
var enableLeaderElection bool
var leaderElectionNamespace string
var metricsPort int
var enableMCP bool
var mcpPort int
var enableOtel bool
var FrontendApi *api.FrontendPageApi
var jwtSecret string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a fasthttp server",
	Run: func(cmd *cobra.Command, args []string) {

		level := SetLogLevel(LogLevel)
		ConfigureLogger(level)

		// Initialize OpenTelemetry if enabled
		var shutdownOtel func(context.Context) error
		if enableOtel {
			config := telemetry.TracingConfig{
				ServiceName:    "k8s-controller-tutorial",
				ServiceVersion: "1.0.0",
				EnableConsole:  true,
			}

			var err error
			shutdownOtel, err = telemetry.InitTracing(context.Background(), config)
			if err != nil {
				log.Error().Err(err).Msg("Failed to initialize OpenTelemetry")
				os.Exit(1)
			}

			// Function for graceful shutdown
			defer func() {
				if shutdownOtel != nil {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					if err := shutdownOtel(ctx); err != nil {
						log.Error().Err(err).Msg("Failed to shutdown OpenTelemetry")
					}
				}
			}()

			log.Info().Msg("OpenTelemetry tracing enabled")
		}

		ctrlruntime.SetLogger(zap.New(zap.UseDevMode(true)))

		clientset, kubeConfig, err := ChooseKubeConnectionType(inCluster, kubeconfig)
		if err != nil {
			log.Error().Err(err).Msg("Error creating clientset")
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go informer.StartDeploymentInformer(ctx, clientset, namespace)

		scheme := runtime.NewScheme()
		if err := clientgoscheme.AddToScheme(scheme); err != nil {
			log.Error().Err(err).Msg("Error adding client-go scheme")
			return
		}
		if err := frontendv1alpha1.AddToScheme(scheme); err != nil {
			log.Error().Err(err).Msg("Error adding frontend scheme")
			os.Exit(1)
		}

		// Start controller-runtime manager and controller
		mgr, err := ctrlruntime.NewManager(kubeConfig, manager.Options{
			Scheme:                  scheme,
			LeaderElection:          enableLeaderElection,
			LeaderElectionID:        "k8s-controller-tutorial-leader",
			LeaderElectionNamespace: leaderElectionNamespace,
			LeaseDuration:           &[]time.Duration{15 * time.Second}[0],
			RenewDeadline:           &[]time.Duration{10 * time.Second}[0],
			RetryPeriod:             &[]time.Duration{2 * time.Second}[0],
			Metrics:                 server.Options{BindAddress: fmt.Sprintf(":%d", metricsPort)},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create controller-runtime manager")
			os.Exit(1)
		}
		if err := ctrl.AddFrontendPageController(mgr); err != nil {
			log.Error().Err(err).Msg("Failed to add frontend page controller")
			os.Exit(1)
		}
		if err := ctrl.AddDeploymentController(mgr); err != nil {
			log.Error().Err(err).Msg("Failed to add deployment controller")
			os.Exit(1)
		}

		router := fasthttprouter.New()

		// Function to wrap handlers with OpenTelemetry middleware
		wrapHandler := func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
			if enableOtel {
				return api.OtelMiddleware(handler)
			}
			return handler
		}

		// Wrap all API endpoints with OpenTelemetry middleware
		router.POST("/api/token", wrapHandler(api.TraceableHandler("GenerateToken", api.TokenHandler)))

		frontedApi := &api.FrontendPageApi{
			K8SClient: mgr.GetClient(),
			Namespace: namespace,
		}
		api.FrontendApi = frontedApi

		router.GET("/api/frontendpages", wrapHandler(api.TraceableHandler("ListFrontendPages", api.JwtMiddleware(frontedApi.ListFrontendPages))))
		router.POST("/api/frontendpages", wrapHandler(api.TraceableHandler("CreateFrontendPage", api.JwtMiddleware(frontedApi.CreateFrontendPage))))
		router.GET("/api/frontendpages/:name", wrapHandler(api.TraceableHandler("GetFrontendPage", api.JwtMiddleware(frontedApi.GetFrontendPage))))
		router.PUT("/api/frontendpages/:name", wrapHandler(api.TraceableHandler("UpdateFrontendPage", api.JwtMiddleware(frontedApi.UpdateFrontendPage))))
		router.DELETE("/api/frontendpages/:name", wrapHandler(api.TraceableHandler("DeleteFrontendPage", api.JwtMiddleware(frontedApi.DeleteFrontendPage))))

		router.GET("/health", wrapHandler(api.TraceableHandler("HealthCheck", func(ctx *fasthttp.RequestCtx) {
			ctx.Response.Header.Set("Content-Type", "application/json")
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.WriteString(`{"status": "ok"}`)
		})))

		CORS := func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
			return func(ctx *fasthttp.RequestCtx) {
				ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
				ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
				ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
				if string(ctx.Method()) == fasthttp.MethodOptions {
					ctx.SetStatusCode(fasthttp.StatusOK)
					return
				}
				h(ctx)
			}
		}
		router.GET("/swagger/*any", CORS(fasthttpadaptor.NewFastHTTPHandler(httpSwagger.WrapHandler)))

		router.GET("/deployments", wrapHandler(api.TraceableHandler("ListDeployments", func(ctx *fasthttp.RequestCtx) {
			ctx.Response.Header.Set("Content-Type", "application/json")
			deployments := informer.GetDeploymentsNames()
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.WriteString("[")
			for i, name := range deployments {
				ctx.WriteString("\"")
				ctx.WriteString(string(name))
				ctx.WriteString("\"")
				if i < len(deployments)-1 {
					ctx.WriteString(",")
				}
			}
			ctx.WriteString("]")
		})))
		api.JWTSecret = jwtSecret

		go func() {
			log.Info().Msg("Starting controller-runtime manager...")
			if err := mgr.Start(cmd.Context()); err != nil {
				log.Error().Err(err).Msg("Manager exited with error")
				os.Exit(1)
			}
		}()

		if enableMCP {
			go func() {
				mcpServer := NewMCPServer("K8S controller MCP", "1.0.0")
				sseServer := mcpserver.NewSSEServer(mcpServer,
					mcpserver.WithBaseURL(fmt.Sprintf("http://:%d/mcp", mcpPort)),
				)
				log.Info().Msgf("Starting MCP server on port %d", mcpPort)
				if err := sseServer.Start(fmt.Sprintf(":%d", mcpPort)); err != nil {
					log.Error().Err(err).Msg("Failed to start SSe server")
					os.Exit(1)
				}
			}()
			log.Info().Msgf("MCP server started on port %d", mcpPort)
		}

		addr := fmt.Sprintf(":%d", serverPort)
		log.Info().Msgf("Starting server on %s", addr)
		if enableOtel {
			log.Info().Msg("OpenTelemetry tracing is enabled for all API endpoints")
		}
		if err := fasthttp.ListenAndServe(addr, router.Handler); err != nil {
			log.Error().Err(err).Msg("Failed to start server")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVar(&serverPort, "port", 8080, "Port to listen on")
	serverCmd.Flags().BoolVar(&inCluster, "in-cluster", false, "Use in-cluster configuration")
	serverCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file")
	serverCmd.Flags().StringVar(&namespace, "namespace", "default", "Namespace to watch")
	serverCmd.Flags().BoolVar(&enableLeaderElection, "leader-election", true, "Enable leader election")
	serverCmd.Flags().StringVar(&leaderElectionNamespace, "leader-election-namespace", "default", "Namespace for leader election")
	serverCmd.Flags().IntVar(&metricsPort, "metrics-port", 8081, "Port for metrics")
	serverCmd.Flags().BoolVar(&enableMCP, "enable-mcp", false, "Enable MCP server")
	serverCmd.Flags().IntVar(&mcpPort, "mcp-port", 9090, "Port for MCP server")
	serverCmd.Flags().BoolVar(&enableOtel, "enable-otel", false, "Enable OpenTelemetry tracing")
	serverCmd.Flags().StringVar(&jwtSecret, "jwt-secret", "secret", "JWT secret (required for token-based authentication)")
}
