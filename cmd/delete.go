package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		level := SetLogLevel(LogLevel)
		ConfigureLogger(level)

		clientset, err := GetKubeClient(kubeconfig)
		if err != nil {
			log.Error().Err(err).Msg("Error creating clientset with kubeconfig file with path: " + kubeconfig)
		}

		log.Info().Msgf("Deleting deployment %s in namespace %s", deploymentName, namespace)
		err = deleteDeployment(clientset, namespace, deploymentName)
		if err != nil {
			log.Error().Err(err).Msg("Error deleting deployment")
			return
		}
		log.Info().Msgf("Deployment %s in namespace %s deleted", deploymentName, namespace)
	},
}

func init() {
	deleteCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file")
	deleteCmd.Flags().StringVar(&namespace, "namespace", "default", "Namespace to delete resource in")
	deleteCmd.Flags().StringVar(&deploymentName, "deployment-name", "my-deployment", "Name of the deployment to delete")
	rootCmd.AddCommand(deleteCmd)
}