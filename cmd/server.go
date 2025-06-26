package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/JRaver/k8s-controller-tutorial/pkg/ctrl"
	"github.com/JRaver/k8s-controller-tutorial/pkg/informer"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	ctrlruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var serverPort int = 8080
var enableLeaderElection bool
var leaderElectionNamespace string
var metricsPort int

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a fasthttp server",
	Run: func(cmd *cobra.Command, args []string) {

		level := SetLogLevel(LogLevel)
		ConfigureLogger(level)

		ctrlruntime.SetLogger(zap.New(zap.UseDevMode(true)))
		clientset, err := ChooseKubeConnectionType(inCluster, kubeconfig)
		if err != nil {
			log.Error().Err(err).Msg("Error creating clientset")
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go informer.StartDeploymentInformer(ctx, clientset, namespace)

		// Start controller-runtime manager and controller
		mgr, err := ctrlruntime.NewManager(ctrlruntime.GetConfigOrDie(), manager.Options{
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
		if err := ctrl.AddDeploymentController(mgr); err != nil {
			log.Error().Err(err).Msg("Failed to add deployment controller")
			os.Exit(1)
		}
		go func() {
			log.Info().Msg("Starting controller-runtime manager...")
			if err := mgr.Start(cmd.Context()); err != nil {
				log.Error().Err(err).Msg("Manager exited with error")
				os.Exit(1)
			}
		}()

		log.Info().Msgf("Starting server on port %d", serverPort)
		handler := func(ctx *fasthttp.RequestCtx) {
			requestID := uuid.New().String()
			ctx.Response.Header.Set("X-Request-ID", requestID)
			logger := log.With().Str("request_id", requestID).Logger()
			path := string(ctx.Path())
			switch path {
			case "/healthz":
				logger.Info().Msg("Health check request")
				ctx.WriteString("OK")
			case "/deployments":
				logger.Info().Msg("Deployments request received")
				ctx.Response.Header.Set("Content-Type", "application/json")
				deployments := informer.GetDeploymentsNames()
				logger.Info().Msgf("Found %d deployments with namespace %s", len(deployments), namespace)
				ctx.SetStatusCode(fasthttp.StatusOK)
				ctx.Write([]byte("["))
				for i, name := range deployments {
					ctx.WriteString("\"")
					ctx.WriteString(string(name))
					ctx.WriteString("\"")
					if i < len(deployments)-1 {
						ctx.WriteString(",")
					}
				}
				ctx.WriteString("]")
				return
			default:
				logger.Info().Msg("Hello from FastHTTP!")
				ctx.WriteString("Hello from FastHTTP!")
			}
		}
		if err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", serverPort), handler); err != nil {
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
}
