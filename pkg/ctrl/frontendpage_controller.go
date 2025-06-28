package ctrl

import (	
	context "context"
	"reflect"

	"github.com/rs/zerolog/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	 frontendv1alpha1 "github.com/JRaver/k8s-controller-tutorial/pkg/apis/frontend/v1alpha1"

)

type FrontendPageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func buildConfigMap(frontendPage *frontendv1alpha1.FrontendPage) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      frontendPage.Name,
			Namespace: frontendPage.Namespace,
		},
		Data: map[string]string{
			"content": frontendPage.Spec.Content,
		},
	}
}

func buildService(frontendPage *frontendv1alpha1.FrontendPage) *corev1.Service {
	port := int32(frontendPage.Spec.Port)
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      frontendPage.Name,
			Namespace: frontendPage.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": frontendPage.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: port,
					TargetPort: intstr.FromInt(int(port)),
				},
			},
		},
	}
}

func buildDeployment(frontendPage *frontendv1alpha1.FrontendPage) *appsv1.Deployment {
	replicas := int32(frontendPage.Spec.Replicas)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      frontendPage.Name,
			Namespace: frontendPage.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": frontendPage.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": frontendPage.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  frontendPage.Name,
							Image: frontendPage.Spec.Image,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "content",
									MountPath: "/data",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "content",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: frontendPage.Name,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
func (r *FrontendPageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var frontendPage frontendv1alpha1.FrontendPage
	err := r.Get(ctx, req.NamespacedName, &frontendPage)
	if err != nil {
		return ctrl.Result{}, err
	}
	
	svc := buildService(&frontendPage)
	if err := ctrl.SetControllerReference(&frontendPage, svc, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	log.Info().Msgf("Reconciled FrontendPage Service: %s/%s", req.Namespace, req.Name)
	
	var existingService corev1.Service
	if err := r.Get(ctx, req.NamespacedName, &existingService); err != nil {
		if !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, svc); err != nil {
			return ctrl.Result{}, err
		}
		log.Info().Msgf("Created FrontendPage Service: %s/%s", req.Namespace, req.Name)
	} else if !reflect.DeepEqual(existingService.Spec, svc.Spec) {
		existingService.Spec = svc.Spec
		if err := r.Update(ctx, &existingService); err != nil {
			return ctrl.Result{}, err
		}
		log.Info().Msgf("Updated FrontendPage Service: %s/%s", req.Namespace, req.Name)
	} else {
		log.Info().Msgf("FrontendPage Service: %s/%s is up to date", req.Namespace, req.Name)
	}
	
	cm := buildConfigMap(&frontendPage)
	if err := ctrl.SetControllerReference(&frontendPage, cm, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	log.Info().Msgf("Reconciled FrontendPage ConfigMap: %s/%s", req.Namespace, req.Name)
	var existingConfigMap corev1.ConfigMap
	
	if err := r.Get(ctx, req.NamespacedName, &existingConfigMap); err != nil {
		if !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}

		if err := r.Create(ctx, cm); err != nil {
			return ctrl.Result{}, err
		}
	} else if !reflect.DeepEqual(existingConfigMap.Data, cm.Data) {
		existingConfigMap.Data = cm.Data
		if err := r.Update(ctx, &existingConfigMap); err != nil {
			return ctrl.Result{}, err
		}
		log.Info().Msgf("Updated FrontendPage ConfigMap: %s/%s", req.Namespace, req.Name)
		return ctrl.Result{Requeue: true}, nil
	} else {
		log.Info().Msgf("FrontendPage ConfigMap: %s/%s is up to date", req.Namespace, req.Name)
		return ctrl.Result{}, nil
	}

	deployment := buildDeployment(&frontendPage)

	if err := ctrl.SetControllerReference(&frontendPage, deployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	log.Info().Msgf("Reconciled FrontendPage Deployment: %s/%s", req.Namespace, req.Name)
	var existingDeployment appsv1.Deployment
	if err := r.Get(ctx, req.NamespacedName, &existingDeployment); err != nil {
		if !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, deployment); err != nil {
			return ctrl.Result{}, err
		}
		log.Info().Msgf("Created FrontendPage Deployment: %s/%s", req.Namespace, req.Name)
	} else {
		updated := false
		if *existingDeployment.Spec.Replicas != *deployment.Spec.Replicas {
			existingDeployment.Spec.Replicas = deployment.Spec.Replicas
			updated = true
		}
		if existingDeployment.Spec.Template.Spec.Containers[0].Image != deployment.Spec.Template.Spec.Containers[0].Image {
			existingDeployment.Spec.Template.Spec.Containers[0].Image = deployment.Spec.Template.Spec.Containers[0].Image
			updated = true
		}
		if updated {
			if err := r.Update(ctx, &existingDeployment); err != nil {
				if !errors.IsConflict(err) {
					return ctrl.Result{}, err
				}
				log.Info().Msgf("Conflict updating FrontendPage Deployment: %s/%s, requeuing", req.Namespace, req.Name)
				return ctrl.Result{Requeue: true}, nil
			}
			log.Info().Msgf("Updated FrontendPage Deployment: %s/%s", req.Namespace, req.Name)
		} else {
			log.Info().Msgf("FrontendPage Deployment: %s/%s is up to date", req.Namespace, req.Name)
		}
	}
	return ctrl.Result{}, nil
}

func AddFrontendPageController(mgr manager.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&frontendv1alpha1.FrontendPage{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&appsv1.Deployment{}).
		Complete(&FrontendPageReconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		})
}