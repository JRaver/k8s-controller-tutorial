// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/frontendpages": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get a list of all frontend pages",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "frontendpages"
                ],
                "summary": "List all frontend pages",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Namespace to filter by",
                        "name": "namespace",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.FrontendPageDocList"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create a new frontend page",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "frontendpages"
                ],
                "summary": "Create a frontend page",
                "parameters": [
                    {
                        "description": "Frontend page to create",
                        "name": "frontendpage",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.FrontendPageDoc"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.FrontendPageDoc"
                        }
                    }
                }
            }
        },
        "/api/frontendpages/{name}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get a frontend page by name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "frontendpages"
                ],
                "summary": "Get a frontend page",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the frontend page",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.FrontendPageDoc"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update a frontend page by name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "frontendpages"
                ],
                "summary": "Update a frontend page",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the frontend page",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.FrontendPageDoc"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Delete a frontend page by name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "frontendpages"
                ],
                "summary": "Delete a frontend page",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the frontend page",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.FrontendPageDoc"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.FrontendPageDoc": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "port": {
                    "type": "integer"
                },
                "replicas": {
                    "type": "integer"
                }
            }
        },
        "api.FrontendPageDocList": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.FrontendPageDoc"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "K8s Controller Tutorial API",
	Description:      "My awesome lab controller with Swagger UI",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
