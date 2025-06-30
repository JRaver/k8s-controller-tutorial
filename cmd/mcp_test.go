package cmd

import (
	"context"
	"testing"

	"github.com/JRaver/k8s-controller-tutorial/pkg/api"
	frontendv1alpha1 "github.com/JRaver/k8s-controller-tutorial/pkg/apis/frontend/v1alpha1"
	"github.com/JRaver/k8s-controller-tutorial/pkg/ctrl"
	"github.com/JRaver/k8s-controller-tutorial/pkg/testutil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// setupTestAPIWithManager is a self-contained helper for MCP integration tests.
func setupTestAPIWithManager(t *testing.T) (*api.FrontendPageApi, client.Client, func()) {
	mgr, k8sClient, _, cleanup := testutil.StartTestManager(t)

	require.NoError(t, ctrl.AddFrontendPageController(mgr))

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_ = mgr.Start(ctx)
	}()

	// Wait for the cache to sync before returning
	if ok := mgr.GetCache().WaitForCacheSync(ctx); !ok {
		cancel()
		t.Fatal("cache did not sync")
	}

	apiInst := &api.FrontendPageApi{
		K8SClient: k8sClient,
		Namespace: "default",
	}
	return apiInst, k8sClient, func() {
		cancel()
		cleanup()
	}
}

func TestMCP_ListFrontendPagesHandler(t *testing.T) {
	apiInst, k8sClient, cleanup := setupTestAPIWithManager(t)
	defer cleanup()
	api.FrontendApi = apiInst

	// Create some FrontendPage resources
	page1 := &frontendv1alpha1.FrontendPage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mcp-page1",
			Namespace: "default",
		},
		Spec: frontendv1alpha1.FrontendPageSpec{
			Content:  "<h1>MCP Page 1</h1>",
			Image:    "nginx:1.21",
			Replicas: 1,
			Port:     8080,
		},
	}
	page2 := &frontendv1alpha1.FrontendPage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mcp-page2",
			Namespace: "default",
		},
		Spec: frontendv1alpha1.FrontendPageSpec{
			Content:  "<h1>MCP Page 2</h1>",
			Image:    "nginx:1.22",
			Replicas: 2,
			Port:     8080,
		},
	}
	require.NoError(t, k8sClient.Create(context.Background(), page1))
	require.NoError(t, k8sClient.Create(context.Background(), page2))

	// Call the shared API logic directly (since MCP handler is not accessible)
	docs, err := api.FrontendApi.ListFrontendPagesRaw(context.Background())
	require.NoError(t, err)
	require.Len(t, docs, 2)
	names := []string{docs[0].Name, docs[1].Name}
	require.Contains(t, names, "mcp-page1")
	require.Contains(t, names, "mcp-page2")
}

func TestMCP_CreateFrontendPageHandler(t *testing.T) {
	apiInst, k8sClient, cleanup := setupTestAPIWithManager(t)
	defer cleanup()
	api.FrontendApi = apiInst

	// Create MCP request
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "create_frontendpage",
			Arguments: map[string]interface{}{
				"name":     "test-create-page",
				"contents": "<h1>Test Create Page</h1>",
				"image":    "nginx:alpine",
				"replicas": 3,
				"port":     9090,
			},
		},
	}

	// Call the handler
	result, err := createFrontendPageHandler(context.Background(), req)
	require.NoError(t, err)
	require.False(t, result.IsError)

	// Verify the page was created
	var page frontendv1alpha1.FrontendPage
	err = k8sClient.Get(context.Background(), client.ObjectKey{
		Namespace: "default",
		Name:      "test-create-page",
	}, &page)
	require.NoError(t, err)
	require.Equal(t, "test-create-page", page.Name)
	require.Equal(t, "<h1>Test Create Page</h1>", page.Spec.Content)
	require.Equal(t, "nginx:alpine", page.Spec.Image)
	require.Equal(t, 3, page.Spec.Replicas)
	require.Equal(t, 9090, page.Spec.Port)
}
