# Kubernetes Controller Tutorial

[![CI](https://github.com/JRaver/k8s-controller-tutorial/actions/workflows/ci.yaml/badge.svg)](https://github.com/JRaver/k8s-controller-tutorial/actions/workflows/ci.yaml)

A comprehensive Kubernetes controller implementation demonstrating best practices for building custom controllers with Go. This project showcases a deployment watcher that automatically creates and manages ConfigMaps based on deployment lifecycle events.

## 🚀 Features

- **Deployment Watcher**: Monitors Kubernetes deployments in real-time using informers
- **Automatic ConfigMap Management**: Creates/updates/deletes ConfigMaps based on deployment events
- **HTTP Server**: FastHTTP-based server with health checks and metrics endpoints
- **CLI Interface**: Cobra-powered command-line interface with multiple subcommands
- **Comprehensive Testing**: Unit tests with envtest framework
- **Helm Charts**: Ready-to-deploy Helm charts for Kubernetes
- **Docker Support**: Multi-stage Docker build with distroless base image
- **Structured Logging**: ZeroLog-based logging with configurable levels
- **OpenTelemetry**: Distributed tracing support

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
  --log-level debug

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

### API Endpoints

The server exposes the following HTTP endpoints:

- `GET /healthz` - Health check endpoint
- `GET /metrics` - Metrics endpoint (placeholder)
- `GET /deployments` - List all deployments in the watched namespace
- `GET /` - Default welcome message

## 🏗️ Architecture

### Core Components

1. **Informer Package** (`pkg/informer/`)
   - Deployment watcher using Kubernetes informers
   - ConfigMap lifecycle management
   - Event-driven architecture

2. **CLI Commands** (`cmd/`)
   - Server command for HTTP server
   - Create/Delete/List commands for resource management
   - Root command with logging configuration

3. **Kubernetes Integration**
   - Client-go for Kubernetes API access
   - Controller-runtime for controller patterns
   - Support for both in-cluster and external configurations

### Controller Logic

The controller implements the following workflow:

1. **Deployment Added**: Creates a new ConfigMap with deployment metadata
2. **Deployment Updated**: Updates existing ConfigMap with incremented counter
3. **Deployment Deleted**: Removes associated ConfigMap

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test ./pkg/informer/...
```

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
│   ├── server.go          # HTTP server implementation
│   ├── create.go          # Resource creation commands
│   ├── delete.go          # Resource deletion commands
│   └── list.go            # Resource listing commands
├── pkg/                   # Core packages
│   ├── informer/          # Kubernetes informer implementation
│   └── testutil/          # Testing utilities
├── charts/                # Helm charts
├── main.go               # Application entry point
├── Dockerfile            # Multi-stage Docker build
├── makefile              # Build and development tasks
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