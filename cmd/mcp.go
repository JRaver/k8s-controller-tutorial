package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/JRaver/k8s-controller-tutorial/pkg/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// NewMCPServer creates and configures a new MCP server for FrontendPage tools
func NewMCPServer(serverName, version string) *server.MCPServer {
	s := server.NewMCPServer(
		serverName,
		version,
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// List tool
	listTool := mcp.NewTool("list_frontendpages",
		mcp.WithDescription("List all FrontendPage resources"),
	)
	// Create tool
	createTool := mcp.NewTool("create_frontendpage",
		mcp.WithDescription("Create a new FrontendPage resource"),
		mcp.WithString("name", mcp.Description("Name of the FrontendPage")),
		mcp.WithString("contents", mcp.Description("HTML contents")),
		mcp.WithString("image", mcp.Description("Container image")),
		mcp.WithNumber("replicas", mcp.Description("Number of replicas")),
	)
	// Delete tool
	deleteTool := mcp.NewTool("delete_frontendpage",
		mcp.WithDescription("Delete a FrontendPage resource"),
		mcp.WithString("name", mcp.Description("Name of the FrontendPage to delete")),
	)
	// TODO: Add update tools as needed

	s.AddTool(listTool, listFrontendPagesHandler)
	s.AddTool(createTool, createFrontendPageHandler)
	s.AddTool(deleteTool, deleteFrontendPageHandler)
	// TODO: Register update handlers

	return s
}

func listFrontendPagesHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if api.FrontendApi == nil {
		return mcp.NewToolResultText("FrontendPageApi is not initialized"), nil
	}
	docs, err := api.FrontendApi.ListFrontendPagesRaw(ctx)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error listing FrontendPages: %v", err)), nil
	}
	jsonBytes, err := json.MarshalIndent(docs, "", "  ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error marshaling result: %v", err)), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func createFrontendPageHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if api.FrontendApi == nil {
		return mcp.NewToolResultText("FrontendPageApi is not initialized"), nil
	}

	name := req.GetString("name", "")
	contents := req.GetString("contents", "")
	image := req.GetString("image", "")
	replicas := req.GetInt("replicas", 1)
	port := req.GetInt("port", 8080)

	doc := api.FrontendPageDoc{
		Name:     name,
		Content:  contents,
		Image:    image,
		Replicas: replicas,
		Port:     port,
	}

	err := api.FrontendApi.CreateFrontendPageRaw(ctx, doc)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error creating FrontendPage: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("FrontendPage '%s' created successfully", name)), nil
}

func deleteFrontendPageHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if api.FrontendApi == nil {
		return mcp.NewToolResultText("FrontendPageApi is not initialized"), nil
	}

	name := req.GetString("name", "")

	err := api.FrontendApi.DeleteFrontendPageRaw(ctx, name)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error deleting FrontendPage: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("FrontendPage '%s' deleted successfully", name)), nil
}
