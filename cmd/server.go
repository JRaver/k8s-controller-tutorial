package cmd

import (
	"os"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
)

var serverPort int = 8080

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a fasthttp server",
	Run: func(cmd *cobra.Command, args []string) {

		level := SetLogLevel(LogLevel)
		ConfigureLogger(level)

		log.Info().Msgf("Starting server on port %d", serverPort)
		handler := func(ctx *fasthttp.RequestCtx) {
			ctx.WriteString("Hello, World!")
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
}