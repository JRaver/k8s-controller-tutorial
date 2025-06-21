# K8s Controller Tutorial Helm Chart

This Helm chart deploys the k8s-controller-tutorial application to Kubernetes.

## Parameters

### Global Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of replicas to deploy | `1` |

### Image Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `image.repository` | Docker image repository | `docker-hub.XXX/k8s-controller-tutorial` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `image.tag` | Image tag override | `""` (uses chart appVersion) |
| `imagePullSecrets` | Image pull secrets | `[]` |

### Naming Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `nameOverride` | Override the chart name | `""` |
| `fullnameOverride` | Override the full chart name | `""` |

### Service Account Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `serviceAccount.create` | Create a service account | `true` |
| `serviceAccount.automount` | Automount API credentials | `true` |
| `serviceAccount.annotations` | Service account annotations | `{}` |
| `serviceAccount.name` | Service account name | `""` (auto-generated) |

### Pod Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `podAnnotations` | Pod annotations | `{}` |
| `podLabels` | Pod labels | `{}` |
| `podSecurityContext` | Pod security context | `{}` |
| `securityContext` | Container security context | `{}` |

### Service Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `service.type` | Service type | `ClusterIP` |
| `service.port` | Service port | `80` |

### Ingress Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `ingress.enabled` | Enable ingress | `false` |
| `ingress.className` | Ingress class name | `""` |
| `ingress.annotations` | Ingress annotations | `{}` |
| `ingress.hosts` | Ingress hosts configuration | `[{"host": "chart-example.local", "paths": [{"path": "/", "pathType": "ImplementationSpecific"}]}]` |
| `ingress.tls` | Ingress TLS configuration | `[]` |

### Resource Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `resources` | CPU/Memory resource requests and limits | `{}` |

### Autoscaling Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `autoscaling.enabled` | Enable horizontal pod autoscaler | `false` |
| `autoscaling.minReplicas` | Minimum number of replicas | `1` |
| `autoscaling.maxReplicas` | Maximum number of replicas | `100` |
| `autoscaling.targetCPUUtilizationPercentage` | Target CPU utilization | `80` |

### Volume Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `volumes` | Additional volumes | `[]` |
| `volumeMounts` | Additional volume mounts | `[]` |

### Scheduling Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `nodeSelector` | Node selector | `{}` |
| `tolerations` | Tolerations | `[]` |
| `affinity` | Affinity rules | `{}` |

## Usage

### Basic Installation

```bash
helm install my-release ./charts/k8s-controller-tutorial
```

### Custom Configuration

```bash
helm install my-release ./charts/k8s-controller-tutorial \
  --set replicaCount=3 \
  --set image.repository=my-registry/k8s-controller-tutorial \
  --set ingress.enabled=true
```

### Using values.yaml

```bash
helm install my-release ./charts/k8s-controller-tutorial -f custom-values.yaml
```

## Example custom-values.yaml

```yaml
replicaCount: 3

image:
  repository: my-registry/k8s-controller-tutorial
  tag: "v1.0.0"

ingress:
  enabled: true
  hosts:
    - host: my-app.example.com
      paths:
        - path: /
          pathType: Prefix

resources:
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 250m
    memory: 256Mi

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
``` 