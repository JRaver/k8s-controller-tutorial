package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		level := SetLogLevel(LogLevel)
		ConfigureLogger(level)

		log.Debug().Msg("Getting kubeclient")
		clientset, err := GetKubeClient(kubeconfig)
		if err != nil {
			log.Error().Err(err).Msg("Error creating clientset with kubeconfig file with path: " + kubeconfig)
			return
		}
		log.Info().Msg("Get existing deployments before creation")
		existingDeployments, err := getExistingDeployments(clientset, namespace)
		if err != nil {
			log.Error().Err(err).Msg("Error getting existing deployments")
			return
		}
		log.Info().Msgf("Found %d deployments", len(existingDeployments.Items))
		for _, deployment := range existingDeployments.Items {
			log.Info().Msgf("Deployment: %s", deployment.Name)
		}

		log.Info().Msgf("Creating deployment %s in namespace %s", deploymentName, namespace)
		deployment := deploymentSpecBuilder(deploymentName, namespace)

		createdDeployment, err := clientset.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
		if err != nil {
			log.Error().Err(err).Msg("Error creating deployment")
			return
		}
		log.Info().Msgf("Successfully created deployment: %s", createdDeployment.Name)

		log.Info().Msg("Get existing deployments after creation")
		existingDeployments, err = getExistingDeployments(clientset, namespace)
		if err != nil {
			log.Error().Err(err).Msg("Error getting existing deployments")
			return
		}
		log.Info().Msgf("Found %d deployments", len(existingDeployments.Items))
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file")
	createCmd.Flags().StringVar(&namespace, "namespace", "default", "Namespace to create resource in")
	createCmd.Flags().StringVar(&deploymentName, "deployment-name", "my-deployment", "Name of the deployment to create")
}
