package ctrl

import (
	"context"

	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type DeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Info().Msgf("Reconciling Deployment: %s/%s", req.Namespace, req.Name)

	// Get Deployment
	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, req.NamespacedName, deployment)
	if err != nil {
		// Deployment not found - maybe deleted
		log.Info().Msgf("Deployment DELETED: %s in namespace %s", req.Name, req.Namespace)
		return ctrl.Result{}, nil
	}

	// Deployment exists - log information
	log.Info().Msgf("Deployment %s in namespace %s exists with %d replicas",
		deployment.Name, deployment.Namespace, *deployment.Spec.Replicas)

	return ctrl.Result{}, nil
}

// Predicate for filtering events
type DeploymentPredicate struct {
	predicate.Funcs
}

// Create - called when Deployment is created
func (DeploymentPredicate) Create(e event.CreateEvent) bool {
	log.Info().Msgf("Deployment CREATED: %s in namespace %s", e.Object.GetName(), e.Object.GetNamespace())
	return true
}

// Update - called when Deployment is updated
func (DeploymentPredicate) Update(e event.UpdateEvent) bool {
	log.Info().Msgf("Deployment UPDATED: %s in namespace %s", e.ObjectNew.GetName(), e.ObjectNew.GetNamespace())
	return true
}

// Delete - called when Deployment is deleted
func (DeploymentPredicate) Delete(e event.DeleteEvent) bool {
	log.Info().Msgf("Deployment DELETED: %s in namespace %s", e.Object.GetName(), e.Object.GetNamespace())
	return true
}

// Generic - called for generic events
func (DeploymentPredicate) Generic(e event.GenericEvent) bool {
	log.Info().Msgf("Deployment GENERIC EVENT: %s in namespace %s", e.Object.GetName(), e.Object.GetNamespace())
	return true
}

func AddDeploymentController(mgr manager.Manager) error {
	r := &DeploymentReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		WithEventFilter(DeploymentPredicate{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}
