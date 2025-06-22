package cmd

import (
	"os"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"context"
	"github.com/JRaver/k8s-controller-tutorial/pkg/informer"
	"github.com/google/uuid"
)

var serverPort int = 8080

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a fasthttp server",
	Run: func(cmd *cobra.Command, args []string) {

		level := SetLogLevel(LogLevel)
		ConfigureLogger(level)

		clientset, err := ChooseKubeConnectionType(inCluster, kubeconfig)
		if err != nil {
			log.Error().Err(err).Msg("Error creating clientset")
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go informer.StartDeploymentInformer(ctx, clientset, namespace)

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
			case "/metrics":
				logger.Info().Msg("Metrics request")
				fmt.Println("Metrics will be realized later")
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
}