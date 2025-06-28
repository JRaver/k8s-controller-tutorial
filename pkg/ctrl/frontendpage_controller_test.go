package ctrl

import (
	context "context"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	frontendv1alpha1 "github.com/JRaver/k8s-controller-tutorial/pkg/apis/frontend/v1alpha1"
	testutil "github.com/JRaver/k8s-controller-tutorial/pkg/testutil"
	"github.com/stretchr/testify/require"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func printTableState(ctx context.Context, c client.Client, ns string, t *testing.T, step string) {
	var pages frontendv1alpha1.FrontendPageList
	var cms corev1.ConfigMapList
	var deps appsv1.DeploymentList
	var svcs corev1.ServiceList

	c.List(ctx, &pages, client.InNamespace(ns))
	c.List(ctx, &cms, client.InNamespace(ns))
	c.List(ctx, &deps, client.InNamespace(ns))

	t.Logf("\n===== ETCD STATE (%s) =====", step)
	t.Logf("%-15s %-15s %-10s %-10s", "KIND", "NAME", "NAMESPACE", "EXTRA")
	for _, p := range pages.Items {
		t.Logf("%-15s %-15s %-10s %-10s", "FrontendPage", p.Name, p.Namespace, p.Spec.Content)
	}
	for _, c := range cms.Items {
		t.Logf("%-15s %-15s %-10s %-10s", "ConfigMap", c.Name, c.Namespace, c.Data["content"])
	}
	for _, s := range svcs.Items {
		t.Logf("%-15s %-15s %-10s %-10s", "Service", s.Name, s.Namespace, s.Spec.Selector["app"])
	}

	for _, d := range deps.Items {
		replicas := int32(0)
		if d.Spec.Replicas != nil {
			replicas = *d.Spec.Replicas
		}
		t.Logf("%-15s %-15s %-10s replicas=%d", "Deployment", d.Name, d.Namespace, replicas)
	}
	if len(pages.Items) == 0 && len(cms.Items) == 0 && len(deps.Items) == 0 {
		t.Logf("<empty>")
	}
}

func TestFrontendPageReconciler_Reconcile(t *testing.T) {
	log.SetLogger(zap.New(zap.UseDevMode(true)))

	_, k8sClient, restCfg, cleanup := testutil.StartTestManager(t)
	defer cleanup()

	ctx := context.Background()
	ns := "default"

	// Check if CRD is installed
	extClient, err := apiextensionsclient.NewForConfig(restCfg)
	require.NoError(t, err)

	crd, err := extClient.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, "frontendpages.frontend.jraver.io", metav1.GetOptions{})
	require.NoError(t, err)
	require.NotNil(t, crd)
	require.Equal(t, "frontendpages.frontend.jraver.io", crd.Name)

	printTableState(ctx, k8sClient, ns, t, "Before Create")

	page := &frontendv1alpha1.FrontendPage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testpage",
			Namespace: ns,
		},
		Spec: frontendv1alpha1.FrontendPageSpec{
			Image:    "nginx:latest",
			Content:  "this is a test",
			Replicas: 1,
			Port:     8888,
		},
	}

	if err := k8sClient.Create(ctx, page); err != nil {
		t.Fatalf("Failed to create FrontendPage: %v", err)
	}

	time.Sleep(1 * time.Second)
	printTableState(ctx, k8sClient, ns, t, "After Create")

	var pageList frontendv1alpha1.FrontendPageList
	err = k8sClient.List(ctx, &pageList, client.InNamespace(ns))
	require.NoError(t, err)
	require.NotEmpty(t, pageList.Items)
	require.Len(t, pageList.Items, 1)
	found := false
	for _, p := range pageList.Items {
		if p.Name == "testpage" && p.Spec.Content == "this is a test" && p.Spec.Port == 8888 {
			found = true
			break
		}
	}

	require.True(t, found, "Created FrontendPage should be found")
	//Update the page
	page.Spec.Content = "Updated Content"

	if err := k8sClient.Update(ctx, page); err != nil {
		t.Fatalf("Failed to update FrontendPage: %v", err)
	}

	time.Sleep(1 * time.Second)
	printTableState(ctx, k8sClient, ns, t, "Updated Content")

	//Delete the page
	if err := k8sClient.Delete(ctx, page); err != nil {
		t.Fatalf("Failed to delete FrontendPage: %v", err)
	}

	time.Sleep(1 * time.Second)
	printTableState(ctx, k8sClient, ns, t, "After Delete")

}
