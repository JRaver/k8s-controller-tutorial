#!/bin/bash

# OpenTelemetry Demo Script for k8s-controller-tutorial
# This script demonstrates how to use OpenTelemetry tracing with the API

set -e

echo "=== OpenTelemetry Demo for k8s-controller-tutorial ==="
echo ""

# Check if binary exists
if [ ! -f "./bin/k8s-controller-tutorial" ]; then
    echo "Building the application..."
    go build -o bin/k8s-controller-tutorial main.go
fi

# Start the server with OpenTelemetry enabled in background
echo "Starting server with OpenTelemetry enabled..."
./bin/k8s-controller-tutorial server \
    --enable-otel \
    --log-level=info \
    --port=8080 \
    --namespace=default \
    --jwt-secret=demo-secret &

SERVER_PID=$!

# Wait for server to start
echo "Waiting for server to start..."
sleep 5

# Check if server is running
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "Server failed to start"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi

echo "Server started successfully!"
echo ""

# Generate JWT token
echo "=== Generating JWT Token ==="
TOKEN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/token \
    -H "Content-Type: application/json" \
    -d '{"username": "admin", "password": "secret"}')

TOKEN=$(echo $TOKEN_RESPONSE | jq -r '.token')
echo "Token generated: ${TOKEN:0:20}..."
echo ""

# Test API endpoints with tracing
echo "=== Testing API Endpoints (watch logs for OpenTelemetry spans) ==="
echo ""

echo "1. Health check (should show tracing in logs)"
curl -s http://localhost:8080/health | jq .
echo ""

echo "2. List frontend pages (should show K8s operation tracing)"
curl -s -X GET http://localhost:8080/api/frontendpages \
    -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "3. Create frontend page (should show detailed K8s create operation tracing)"
curl -s -X POST http://localhost:8080/api/frontendpages \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "demo-page",
        "content": "Hello OpenTelemetry!",
        "image": "nginx:alpine",
        "replicas": 1,
        "port": 80
    }' | jq .
echo ""

echo "4. Get specific frontend page (should show K8s get operation tracing)"
curl -s -X GET http://localhost:8080/api/frontendpages/demo-page \
    -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "5. Update frontend page (should show K8s update operation tracing)"
curl -s -X PUT http://localhost:8080/api/frontendpages/demo-page \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "demo-page",
        "content": "Updated with OpenTelemetry!",
        "image": "nginx:alpine",
        "replicas": 2,
        "port": 80
    }' | jq .
echo ""

echo "6. Delete frontend page (should show K8s delete operation tracing)"
curl -s -X DELETE http://localhost:8080/api/frontendpages/demo-page \
    -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "=== Demo Complete ==="
echo ""
echo "Check the server logs above to see OpenTelemetry traces with:"
echo "- HTTP request spans with method, path, status code, duration"
echo "- Kubernetes operation spans with namespace, resource names, attributes"
echo "- Span events for operation start/end"
echo "- Error recording for failed operations"
echo ""
echo "Press Ctrl+C to stop the server"

# Wait for user to stop
trap 'kill $SERVER_PID 2>/dev/null || true; exit 0' INT
wait $SERVER_PID 