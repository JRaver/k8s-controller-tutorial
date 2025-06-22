package informer

import (
	"context"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var deploymentInformer cache.SharedIndexInformer

func StartDeploymentInformer(ctx context.Context, clientset *kubernetes.Clientset, namespace string) {
	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset,
		10*time.Second,
		informers.WithNamespace(namespace),
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.FieldSelector = fields.Everything().String()
		}),
	)
	deploymentInformer = factory.Apps().V1().Deployments().Informer()
	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Info().Msg("Deployment added: " + GetDeploymentName(obj))
			configMap := ConfigMapBuilder(GetDeploymentName(obj), namespace, GetDeploymentName(obj), 0)
			clientset.CoreV1().ConfigMaps(namespace).Create(context.Background(), configMap, metav1.CreateOptions{})
			log.Info().Msg("ConfigMap created: " + GetDeploymentName(obj))
		},

		UpdateFunc: func(oldObj, newObj interface{}) {
			//log.Info().Msg("Deployment updated: " + GetDeploymentName(newObj))
			oldDeployment := oldObj.(metav1.Object)
			newDeployment := newObj.(metav1.Object)

			if oldDeployment.GetResourceVersion() != newDeployment.GetResourceVersion() {
				log.Info().Msg("Deployment updated: " + newDeployment.GetName())

				if configMapExists, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), GetDeploymentName(newObj), metav1.GetOptions{}); err != nil {
					configMap := ConfigMapBuilder(GetDeploymentName(newObj), namespace, GetDeploymentName(newObj), 0)
					clientset.CoreV1().ConfigMaps(namespace).Create(context.Background(), configMap, metav1.CreateOptions{})

					log.Info().Msg("ConfigMap created: " + GetDeploymentName(newObj))
				} else {
					counter, _ := strconv.Atoi(configMapExists.Data["updateCount"])
					counter++
					log.Info().Msg("ConfigMap found: " + newDeployment.GetName() + " with updateCount: " + strconv.Itoa(counter))
					configMap := ConfigMapBuilder(GetDeploymentName(newObj), namespace, GetDeploymentName(newObj), counter)
					clientset.CoreV1().ConfigMaps(namespace).Update(context.Background(), configMap, metav1.UpdateOptions{})
					log.Info().Msg("ConfigMap updated: " + newDeployment.GetName())
				}

			}
		},

		DeleteFunc: func(obj interface{}) {
			log.Info().Msg("Deployment deleted: " + GetDeploymentName(obj))
			clientset.CoreV1().ConfigMaps(namespace).Delete(context.Background(), GetDeploymentName(obj), metav1.DeleteOptions{})
		},
	})
	log.Info().Msg("Starting deployment informer with namespace: " + namespace)
	factory.Start(ctx.Done())
	for _, ok := range factory.WaitForCacheSync(ctx.Done()) {
		if !ok {
			log.Error().Msg("Failed to sync cache")
		}
	}
	log.Info().Msg("Deployment informer started, watching for deployments in namespace: " + namespace)
	<-ctx.Done()
	log.Info().Msg("Deployment informer stopped")
}

func GetDeploymentName(obj any) string {
	deployment, ok := obj.(metav1.Object)
	if !ok {
		log.Error().Msg("Object is not a deployment")
		return "unknown"
	}
	return deployment.GetName()
}

func GetDeploymentsNames() []string {
	var names []string
	if deploymentInformer == nil {
		return names
	}
	for _, obj := range deploymentInformer.GetStore().List() {
		deployment, ok := obj.(metav1.Object)
		if !ok {
			log.Error().Msg("Object is not a deployment")
			continue
		}
		names = append(names, deployment.GetName())
	}
	return names
}

func ConfigMapBuilder(configMapName string, namespace string, deploymentName string, counter int) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: configMapName,
		},
		Data: map[string]string{
			"deploymentName": deploymentName,
			"updatedAt":      time.Now().Format(time.RFC3339),
			"updateCount":    strconv.Itoa(counter),
		},
	}
}
