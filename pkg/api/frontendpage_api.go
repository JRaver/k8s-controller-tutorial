package api

import (
	"context"
	"encoding/json"
	"fmt"

	frontendv1alpha1 "github.com/JRaver/k8s-controller-tutorial/pkg/apis/frontend/v1alpha1"
	"github.com/valyala/fasthttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// FrontendPageApi is the API for the frontend page
type FrontendPageApi struct {
	K8SClient client.Client
	Namespace string
}

// FrontendPageApi is a shared instance for use by HTTP and MCP handlers
var FrontendApi *FrontendPageApi

// --- Swagger-only structs for doc
// FrontendPageDoc is a simplified version of the FrontendPage type for swagger
type FrontendPageDoc struct {
	Name     string `json:"name"`
	Content  string `json:"content"`
	Image    string `json:"image"`
	Replicas int    `json:"replicas"`
	Port     int    `json:"port"`
}

// FrontendPageDocList is a list of FrontendPageDoc
type FrontendPageDocList struct {
	Items []FrontendPageDoc `json:"items"`
}

// --- API methods
// ListFrontendPages godoc
// @Summary List all frontend pages
// @Description Get a list of all frontend pages
// @Tags frontendpages
// @Accept json
// @Produce json
// @Success 200 {object} FrontendPageDocList
// @Router /api/frontendpages [get]
// @Security ApiKeyAuth
// @Param namespace query string false "Namespace to filter by"

func (api *FrontendPageApi) ListFrontendPages(ctx *fasthttp.RequestCtx) {
	docs, err := api.ListFrontendPagesRaw(context.Background())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
		return
	}

	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(FrontendPageDocList{
		Items: docs,
	})
}

// ListFrontendPagesRaw returns raw list of frontend pages
func (api *FrontendPageApi) ListFrontendPagesRaw(ctx context.Context) ([]FrontendPageDoc, error) {
	list := &frontendv1alpha1.FrontendPageList{}
	if err := api.K8SClient.List(ctx, list, client.InNamespace(api.Namespace)); err != nil {
		return nil, err
	}

	docs := make([]FrontendPageDoc, 0, len(list.Items))
	for _, page := range list.Items {
		docs = append(docs, FrontendPageDoc{
			Name:     page.Name,
			Content:  page.Spec.Content,
			Image:    page.Spec.Image,
			Replicas: page.Spec.Replicas,
			Port:     page.Spec.Port,
		})
	}
	return docs, nil
}

// GetFrontendPage godoc
// @Summary Get a frontend page
// @Description Get a frontend page by name
// @Tags frontendpages
// @Accept json
// @Produce json
// @Success 200 {object} FrontendPageDoc
// @Router /api/frontendpages/{name} [get]
// @Security ApiKeyAuth
// @Param name path string true "Name of the frontend page"

func (api *FrontendPageApi) GetFrontendPage(ctx *fasthttp.RequestCtx) {
	nameValue := ctx.UserValue("name")
	if nameValue == nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(`{"error": "name is required"}`)
		return
	}

	name := nameValue.(string)
	page := &frontendv1alpha1.FrontendPage{}
	err := api.K8SClient.Get(context.Background(), client.ObjectKey{Namespace: api.Namespace, Name: name}, page)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
		return
	}

	// Convert to FrontendPageDoc
	doc := FrontendPageDoc{
		Name:     page.Name,
		Content:  page.Spec.Content,
		Image:    page.Spec.Image,
		Replicas: page.Spec.Replicas,
		Port:     page.Spec.Port,
	}

	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(doc)
}

// CreateFrontendPageRaw creates a frontend page directly (for MCP usage)
func (api *FrontendPageApi) CreateFrontendPageRaw(ctx context.Context, doc FrontendPageDoc) error {
	// Validate required fields
	if doc.Name == "" {
		return fmt.Errorf("name is required")
	}

	// Create FrontendPage from FrontendPageDoc
	object := &frontendv1alpha1.FrontendPage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      doc.Name,
			Namespace: api.Namespace,
		},
		Spec: frontendv1alpha1.FrontendPageSpec{
			Content:  doc.Content,
			Image:    doc.Image,
			Replicas: doc.Replicas,
			Port:     doc.Port,
		},
	}

	return api.K8SClient.Create(ctx, object)
}

// CreateFrontendPage godoc
// @Summary Create a frontend page
// @Description Create a new frontend page
// @Tags frontendpages
// @Accept json
// @Produce json
// @Success 200 {object} FrontendPageDoc
// @Router /api/frontendpages [post]
// @Security ApiKeyAuth
// @Param frontendpage body FrontendPageDoc true "Frontend page to create"
func (api *FrontendPageApi) CreateFrontendPage(ctx *fasthttp.RequestCtx) {
	// Parse FrontendPageDoc from request body
	var doc FrontendPageDoc
	if err := json.Unmarshal(ctx.PostBody(), &doc); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
		return
	}

	// Validate required fields
	if doc.Name == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(`{"error": "name is required"}`)
		return
	}

	// Create FrontendPage from FrontendPageDoc
	object := &frontendv1alpha1.FrontendPage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      doc.Name,
			Namespace: api.Namespace,
		},
		Spec: frontendv1alpha1.FrontendPageSpec{
			Content:  doc.Content,
			Image:    doc.Image,
			Replicas: doc.Replicas,
			Port:     doc.Port,
		},
	}

	if err := api.K8SClient.Create(context.Background(), object); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusCreated)
	json.NewEncoder(ctx).Encode(doc)
}

// UpdateFrontendPage godoc
// @Summary Update a frontend page
// @Description Update a frontend page by name
// @Tags frontendpages
// @Accept json
// @Produce json
// @Success 200 {object} FrontendPageDoc
// @Router /api/frontendpages/{name} [put]
// @Security ApiKeyAuth
// @Param name path string true "Name of the frontend page"

func (api *FrontendPageApi) UpdateFrontendPage(ctx *fasthttp.RequestCtx) {
	nameValue := ctx.UserValue("name")
	if nameValue == nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(`{"error": "name is required"}`)
		return
	}

	name := nameValue.(string)

	// Fetch the existing page
	existingPage := &frontendv1alpha1.FrontendPage{}
	err := api.K8SClient.Get(context.Background(), client.ObjectKey{Namespace: api.Namespace, Name: name}, existingPage)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
		return
	}

	// Parse FrontendPageDoc from request body
	var doc FrontendPageDoc
	if err := json.Unmarshal(ctx.PostBody(), &doc); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
		return
	}

	// Update the page spec
	existingPage.Spec.Content = doc.Content
	existingPage.Spec.Image = doc.Image
	existingPage.Spec.Replicas = doc.Replicas
	existingPage.Spec.Port = doc.Port

	if err := api.K8SClient.Update(context.Background(), existingPage); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
		return
	}

	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(doc)
}

// DeleteFrontendPageRaw deletes a frontend page directly (for MCP usage)
func (api *FrontendPageApi) DeleteFrontendPageRaw(ctx context.Context, name string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}

	return api.K8SClient.Delete(ctx, &frontendv1alpha1.FrontendPage{
		ObjectMeta: metav1.ObjectMeta{Namespace: api.Namespace, Name: name},
	})
}

// DeleteFrontendPage godoc
// @Summary Delete a frontend page
// @Description Delete a frontend page by name
// @Tags frontendpages
// @Accept json
// @Produce json
// @Success 200 {object} FrontendPageDoc
// @Router /api/frontendpages/{name} [delete]
// @Security ApiKeyAuth
// @Param name path string true "Name of the frontend page"

func (api *FrontendPageApi) DeleteFrontendPage(ctx *fasthttp.RequestCtx) {
	nameValue := ctx.UserValue("name")
	if nameValue == nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(`{"error": "name is required"}`)
		return
	}

	name := nameValue.(string)

	//Delete the page
	if err := api.K8SClient.Delete(context.Background(), &frontendv1alpha1.FrontendPage{
		ObjectMeta: metav1.ObjectMeta{Namespace: api.Namespace, Name: name},
	}); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
		return
	}
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.WriteString(`{"message": "Frontend page deleted successfully"}`)
}
