# Kubernetes Controller Tutorial

[![CI](https://github.com/JRaver/k8s-controller-tutorial/actions/workflows/ci.yaml/badge.svg)](https://github.com/JRaver/k8s-controller-tutorial/actions/workflows/ci.yaml)

A comprehensive Kubernetes controller implementation demonstrating best practices for building custom controllers with Go. This project showcases a deployment watcher that automatically creates and manages ConfigMaps based on deployment lifecycle events.

## 🚀 Progress/Розробка

### Completed
- ✅ **Step 1**: Golang CLI application using cobra-cli
- ✅ **Step 2**: ZeroLog for log levels - info, debug, trace, warn, error  
- ✅ **Step 3**: pflag with flags for logs level
- ✅ **Step 4**: fasthttp with cobra command "server" and flags for server port
- ✅ **Step 4+**: Add http requests logging
- ✅ **Step 5**: makefile, distroless dockerfile, github workflow and initial tests
- ✅ **Step 6**: k8s.io/client-go to create function to list Kubernetes deployment resources in default namespace. Auth via kubeconfig. add flags for set kubeconfig file. list cli command call function
- ✅ **Step 6+**: Add create/delete command
- ✅ **Step 7**: using k8s.io/client-go create list/watch informer for Kubernetes deployment resources. Auth via kubeconfig or in-cluster auth add flags for in-cluster mode. Informer report events in logs. Add envtest unit tests
- ✅ **Step 7+**: add custom logic function for update/delete events using informers cache search
- ✅ **Step 8**: json api handler to request list deployment resources in informer cache storage
- ✅ **Step 9**: controller-runtime and controller with logic to report in log each event received with informer
- ✅ **Step 10**: controller mgr to control informer and controller. Leader election with lease resource. flag to disable leader election. flag for mgr metrics port
- ✅ **Step 11**: custom crd Frontendpage with additional informer, controller with additional reconciliation logic for custom resource
- ✅ **Step 12**: platform engineering integration based on Port.io actions. API handler for actions to CRUD custom resource
- ✅ **Step 12+**: Add Update action support for IDP and controller
- ✅ **Step 13**: github.com/mark3labs/mcp-go/mcp to create mcp server for api handlers as a mcp tools. flag to specify mcp port
- ✅ **Step 13+**: Add delete/update MCP tool
- ✅ **Step 14**: jwt authentication and authorization for api
- ✅ **Step 15**: basic OpenTelemetry code instrumentation:

### TODO:
- ❌ **Step 3+**: Use Viper to add env vars 
- ❌ **Step 7++**: use config to setup informers start configuration 
- ❌ **Step 15++**: 90% test coverage 
- ❌ **Step 9+**: multi-cluster informers. Dynamically created informers
- ❌ **Step 11++**: add multi-project client configuration for management clusters
- ❌ **Step 12++**: Discord notifications integration
- ❌ **Step 13++**: Add oidc auth to MCP
- ❌ **Step 14+**: add jwt auth for MCP

## 🚀 Features

- **Deployment Watcher**: Monitors Kubernetes deployments in real-time using informers
- **Automatic ConfigMap Management**: Creates/updates/deletes ConfigMaps based on deployment events
- **Custom Resource Controller**: Full CRD implementation with FrontendPage custom resource
- **HTTP Server**: FastHTTP-based server with health checks, metrics, and Swagger API endpoints
- **CLI Interface**: Cobra-powered command-line interface with multiple subcommands (server, list, create, delete, mcp)
- **MCP Server**: Model Context Protocol server for AI assistant integration
- **JWT Authentication**: Secure API access with JWT tokens
- **Platform Engineering Integration**: API handlers for self-service experiences
- **Multi-cluster Support**: Client-go and controller-runtime integration
- **Leader Election**: Built-in leader election for high availability
- **Comprehensive Testing**: Unit tests with envtest framework
- **Helm Charts**: Ready-to-deploy Helm charts for Kubernetes
- **Docker Support**: Multi-stage Docker build with distroless base image
- **Structured Logging**: ZeroLog-based logging with configurable levels (trace, debug, info, warn, error)
- **OpenTelemetry**: Distributed tracing support with spans and metrics

## 📋 Prerequisites

- Go 1.24.2+
- Kubernetes cluster (local or remote)
- kubectl configured
- Docker (optional, for containerized deployment)
- Helm 3.x (optional, for Helm deployment)

## 🛠️ Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/JRaver/k8s-controller-tutorial.git
cd k8s-controller-tutorial

# Build the application
make build

# Or build manually
go build -o k8s-controller-tutorial main.go
```

### Using Docker

```bash
# Build Docker image
make docker-build

# Or build manually
docker build -t k8s-controller-tutorial:latest .
```

## 🚀 Usage

### Basic Commands

```bash
# Start the server with default settings
./k8s-controller-tutorial server

# Start server with custom configuration
./k8s-controller-tutorial server \
  --port 9090 \
  --namespace my-namespace \
  --in-cluster \
  --log-level debug \
  --enable-leader-election \
  --metrics-port 8081 \
  --enable-mcp \
  --mcp-port 8082 \
  --enable-otel \
  --jwt-secret "your-secret-key"

# List deployments in cluster
./k8s-controller-tutorial list \
  --kubeconfig ~/.kube/config \
  --namespace default

# Create a deployment
./k8s-controller-tutorial create \
  --kubeconfig ~/.kube/config \
  --namespace default \
  --deployment-name my-app

# Delete a deployment
./k8s-controller-tutorial delete \
  --kubeconfig ~/.kube/config \
  --namespace default \
  --deployment-name my-app

# List available commands
./k8s-controller-tutorial --help
```

### Configuration Options

| Flag | Description | Default |
|------|-------------|---------|
| `--port` | HTTP server port | 8080 |
| `--namespace` | Kubernetes namespace to watch | default |
| `--in-cluster` | Use in-cluster Kubernetes config | false |
| `--kubeconfig` | Path to kubeconfig file | "" |
| `--log-level` | Logging level (trace, debug, info, warn, error) | info |
| `--enable-leader-election` | Enable leader election for high availability | false |
| `--leader-election-namespace` | Namespace for leader election | default |
| `--metrics-port` | Port for metrics endpoint | 8080 |
| `--enable-mcp` | Enable MCP server | false |
| `--mcp-port` | Port for MCP server | 8080 |
| `--enable-otel` | Enable OpenTelemetry tracing | false |
| `--jwt-secret` | JWT secret key for authentication | "" |
| `--deployment-name` | Name of deployment for create/delete operations | "my-deployment" |

### API Endpoints

The server exposes the following HTTP endpoints:

#### Health & Monitoring
- `GET /health` - Health check endpoint
- `GET /metrics` - Metrics endpoint (controller-runtime metrics)

#### Deployments
- `GET /deployments` - List all deployments in the watched namespace

#### Authentication  
- `POST /api/token` - Generate JWT token for authentication

#### FrontendPage API (Custom Resource)
- `GET /api/frontendpages` - List all FrontendPage resources
- `POST /api/frontendpages` - Create a new FrontendPage resource
- `GET /api/frontendpages/{name}` - Get FrontendPage resource by name
- `PUT /api/frontendpages/{name}` - Update FrontendPage resource
- `DELETE /api/frontendpages/{name}` - Delete FrontendPage resource

#### Documentation
- `GET /swagger/*` - Swagger UI for API documentation

#### MCP Server (if enabled)
- MCP tools for AI assistant integration:
  - `list_frontendpages` - List all FrontendPage resources
  - `create_frontendpage` - Create a new FrontendPage resource  
  - `delete_frontendpage` - Delete a FrontendPage resource

> **Note**: All FrontendPage API endpoints require JWT authentication via Authorization header.

## 🏗️ Architecture

### Core Components

1. **Informer Package** (`pkg/informer/`)
   - Deployment watcher using Kubernetes informers
   - ConfigMap lifecycle management
   - Event-driven architecture with caching

2. **Controller Package** (`pkg/ctrl/`)
   - FrontendPage custom resource controller
   - Deployment controller with reconciliation logic
   - Controller-runtime based controllers

3. **API Package** (`pkg/api/`)
   - FastHTTP-based REST API server
   - JWT authentication middleware
   - OpenTelemetry tracing integration
   - Swagger documentation support

4. **CLI Commands** (`cmd/`)
   - Server command for HTTP server with full configuration
   - Create/Delete/List commands for resource management
   - MCP server integration
   - Root command with advanced logging configuration

5. **Kubernetes Integration**
   - Client-go for Kubernetes API access
   - Controller-runtime for controller patterns
   - Support for both in-cluster and external configurations
   - Leader election for high availability

6. **Telemetry Package** (`pkg/telemetry/`)
   - OpenTelemetry tracing setup
   - Span management and error recording
   - Structured logging integration

### Controller Logic

The system implements multiple controller workflows:

#### Deployment Controller
1. **Deployment Added**: Creates a new ConfigMap with deployment metadata
2. **Deployment Updated**: Updates existing ConfigMap with incremented counter
3. **Deployment Deleted**: Removes associated ConfigMap

#### FrontendPage Controller
1. **Resource Created**: Creates ConfigMap, Service, and Deployment based on spec
2. **Resource Updated**: Updates associated Kubernetes resources
3. **Resource Deleted**: Cleanup through owner references

### Authentication & Authorization

- JWT-based authentication for API endpoints
- Token generation endpoint for clients
- Middleware-based authorization

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run controller tests specifically
make test-controller

# Run specific test
go test ./pkg/informer/...

# Run tests with verbose output
go test -v ./...

# Check test coverage
go test -cover ./...
```

## 🤖 MCP (Model Context Protocol) Integration

The project includes MCP server support for AI assistant integration:

### Available MCP Tools

- `list_frontendpages` - List all FrontendPage resources
- `create_frontendpage` - Create a new FrontendPage resource
- `delete_frontendpage` - Delete a FrontendPage resource

### Running MCP Server

```bash
# Start server with MCP enabled
./k8s-controller-tutorial server --enable-mcp --mcp-port 8082

# Or use environment variables
export ENABLE_MCP=true
export MCP_PORT=8082
./k8s-controller-tutorial server
```

### MCP Client Usage

The MCP server can be integrated with AI assistants like Claude Desktop or other MCP-compatible clients to provide Kubernetes resource management capabilities.

## 📦 Deployment

### Using Helm

```bash
# Install using Helm
helm install k8s-controller-tutorial ./charts/k8s-controller-tutorial

# Upgrade existing installation
helm upgrade k8s-controller-tutorial ./charts/k8s-controller-tutorial

# Uninstall
helm uninstall k8s-controller-tutorial
```

### Using kubectl

```bash
# Apply Kubernetes manifests
kubectl apply -f charts/k8s-controller-tutorial/templates/

# Check deployment status
kubectl get pods -l app=k8s-controller-tutorial
```

## 🔧 Development

### Project Structure

```
k8s-controller-tutorial/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command and logging setup
│   ├── server.go          # HTTP server implementation with full features
│   ├── create.go          # Resource creation commands
│   ├── delete.go          # Resource deletion commands
│   ├── list.go            # Resource listing commands
│   ├── mcp.go             # MCP server tools and handlers
│   └── kuberenets_funcs.go # Kubernetes utility functions
├── pkg/                   # Core packages
│   ├── api/               # HTTP API implementation
│   │   ├── frontendpage_api.go    # FrontendPage CRUD API
│   │   ├── jwt_token.go           # JWT token generation
│   │   ├── jwt_middelware.go      # JWT authentication middleware
│   │   ├── otel_middleware.go     # OpenTelemetry middleware
│   │   └── swagger.go             # Swagger documentation
│   ├── apis/frontend/v1alpha1/    # Custom resource definitions
│   │   ├── resource.go            # FrontendPage resource spec
│   │   ├── groupversion_info.go   # API group version info
│   │   └── zz_generated.deepcopy.go # Generated deepcopy methods
│   ├── ctrl/              # Controller implementations
│   │   ├── frontendpage_controller.go # FrontendPage controller
│   │   └── deployment_controller.go   # Deployment controller
│   ├── informer/          # Kubernetes informer implementation
│   │   └── informer.go    # Deployment informer with caching
│   ├── telemetry/         # OpenTelemetry integration
│   │   └── telemetry.go   # Tracing setup and utilities
│   └── testutil/          # Testing utilities
│       └── envtest.go     # Environment test setup
├── config/crd/            # Custom Resource Definitions
│   └── frontendpage.jraver.io_frontendpages.yaml
├── charts/                # Helm charts
├── docs/                  # Generated documentation
│   ├── swagger.json       # Swagger API specification
│   └── swagger.yaml       # Swagger API specification
├── main.go               # Application entry point
├── Dockerfile            # Multi-stage Docker build with distroless
├── Makefile              # Build and development tasks
└── go.mod                # Go module dependencies
```

### Available Make Targets

```bash
make build          # Build the application
make test           # Run tests
make test-coverage  # Run tests with coverage
make docker-build   # Build Docker image
make clean          # Clean build artifacts
make format         # Format Go code
make lint           # Run linter
```

## 📊 Monitoring

The application provides structured logging with configurable levels:

- **Trace**: Detailed function calls and execution flow
- **Debug**: Debug information for troubleshooting
- **Info**: General application information
- **Warn**: Warning messages
- **Error**: Error conditions

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Kubernetes Client-Go](https://github.com/kubernetes/client-go)
- [Controller Runtime](https://github.com/kubernetes-sigs/controller-runtime)
- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [ZeroLog](https://github.com/rs/zerolog) for structured logging
- [FastHTTP](https://github.com/valyala/fasthttp) for HTTP server

## 📞 Support

For questions and support, please open an issue on GitHub or contact the maintainers.

## OpenTelemetry Integration

The project includes comprehensive OpenTelemetry integration that logs spans to the console and provides distributed tracing for all API endpoints.

### Usage

#### Starting the server with OpenTelemetry enabled

```bash
# Run with OpenTelemetry tracing enabled
./bin/k8s-controller-tutorial server --enable-otel --log-level=debug

# Or with additional options
./bin/k8s-controller-tutorial server \
  --enable-otel \
  --port=8080 \
  --namespace=default \
  --log-level=info \
  --jwt-secret=my-secret-key
```

#### Available flags

- `--enable-otel`: Enable OpenTelemetry tracing (default: false)
- `--port`: Server port (default: 8080)
- `--namespace`: Kubernetes namespace to watch (default: default)
- `--log-level`: Log level (default: info)
- `--jwt-secret`: JWT secret for authentication (default: secret)

#### OpenTelemetry Features

When `--enable-otel` is enabled, the application will:

1. **Initialize OpenTelemetry**: Sets up tracing with console exporter
2. **Trace all API endpoints**: Automatically wraps all HTTP handlers
3. **Log spans to console**: Outputs detailed span information to logs
4. **Track HTTP requests**: Records method, path, status code, duration
5. **Trace Kubernetes operations**: Monitors K8s API calls with detailed attributes

### Example Output

When OpenTelemetry is enabled, you'll see detailed trace information in the logs:

```json
{
  "level": "info",
  "time": "2025-01-27T10:00:00Z",
  "message": "Starting span",
  "span_name": "GET /api/frontendpages",
  "trace_id": "abc123def456",
  "span_id": "789xyz012"
}

{
  "level": "info", 
  "time": "2025-01-27T10:00:00Z",
  "message": "Added span event",
  "span_name": "start_listfrontendpages",
  "trace_id": "abc123def456",
  "span_id": "789xyz012"
}

{
  "level": "info",
  "time": "2025-01-27T10:00:00Z", 
  "message": "Operation duration recorded",
  "operation": "HTTP GET",
  "duration_ms": 45,
  "trace_id": "abc123def456",
  "span_id": "789xyz012"
}
```

### API Endpoints with Tracing

All API endpoints are automatically traced when `--enable-otel` is enabled:

- `POST /api/token` - Generate JWT token
- `GET /api/frontendpages` - List all frontend pages
- `POST /api/frontendpages` - Create a new frontend page
- `GET /api/frontendpages/{name}` - Get specific frontend page
- `PUT /api/frontendpages/{name}` - Update frontend page
- `DELETE /api/frontendpages/{name}` - Delete frontend page
- `GET /health` - Health check endpoint
- `GET /deployments` - List deployments

### Traced Operations

The following operations are traced with detailed attributes:

#### HTTP Request Tracing
- HTTP method and path
- Request URL and user agent
- Remote address
- Response status code and size
- Request duration

#### Kubernetes Operations
- List/Get/Create/Update/Delete operations
- Namespace and resource name
- Operation success/failure
- Resource attributes (image, replicas, etc.)

### Testing the Tracing

1. Start the server with tracing enabled:
```bash
./bin/k8s-controller-tutorial server --enable-otel --log-level=debug
```

2. Generate a JWT token:
```bash
curl -X POST http://localhost:8080/api/token \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "secret"}'
```

3. Make API calls and observe trace logs:
```bash
# Get the token from previous response
TOKEN="your-jwt-token-here"

# List frontend pages (will show tracing in logs)
curl -X GET http://localhost:8080/api/frontendpages \
  -H "Authorization: Bearer $TOKEN"

# Create a frontend page (will show detailed K8s operation tracing)
curl -X POST http://localhost:8080/api/frontendpages \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-page",
    "content": "Hello World",
    "image": "nginx:latest",
    "replicas": 2,
    "port": 80
  }'
```

### Building and Running

```bash
# Build the application
go build -o bin/k8s-controller-tutorial main.go

# Run with OpenTelemetry
./bin/k8s-controller-tutorial server --enable-otel

# View help
./bin/k8s-controller-tutorial server --help
```

## Development

### Dependencies

The project uses the following OpenTelemetry libraries:

- `go.opentelemetry.io/otel` - Core OpenTelemetry API
- `go.opentelemetry.io/otel/sdk` - OpenTelemetry SDK
- `go.opentelemetry.io/otel/exporters/stdout/stdouttrace` - Console trace exporter
- `go.opentelemetry.io/otel/trace` - Tracing API

### Architecture

The OpenTelemetry integration consists of:

1. **Telemetry Package** (`pkg/telemetry/`) - Core tracing functionality
2. **API Middleware** (`pkg/api/otel_middleware.go`) - HTTP request tracing
3. **Server Integration** (`cmd/server.go`) - OpenTelemetry initialization

### Customization

To modify tracing behavior, edit the `TracingConfig` in `cmd/server.go`:

```go
config := telemetry.TracingConfig{
    ServiceName:    "k8s-controller-tutorial",
    ServiceVersion: "1.0.0", 
    EnableConsole:  true,
}
``` 