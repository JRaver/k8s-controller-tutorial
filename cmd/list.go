package cmd

import (

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all resources in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		level := SetLogLevel(LogLevel)
		ConfigureLogger(level)

		clientset, err := GetKubeClient(kubeconfig)
		if err != nil {
			log.Error().Err(err).Msg("Error creating clientset with kubeconfig file with path: " + kubeconfig)
		}
		
		log.Info().Msg("Listing deployments")
		
		existingDeployments, err := getExistingDeployments(clientset, namespace)
		if err != nil {
			log.Error().Err(err).Msg("Error getting existing deployments")
			return
		}	

		log.Info().Msgf("Found %d deployments", len(existingDeployments.Items))
		for _, deployment := range existingDeployments.Items {
			log.Info().Msgf("Deployment: %s", deployment.Name)
		}
	},
}

func init() {
	listCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file")
	listCmd.Flags().StringVar(&namespace, "namespace", "default", "Namespace to list resources in")
	rootCmd.AddCommand(listCmd)
}