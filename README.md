# Ember Reflector Lab

A Kubernetes demonstration of the [Emberstack Reflector](https://github.com/emberstack/kubernetes-reflector) operator that automatically reflects ConfigMaps across namespaces.

## What This Demo Does

This demo showcases how to use the Emberstack Reflector to automatically synchronize configuration data across multiple Kubernetes namespaces. Here's what happens:

1. **Master Configuration**: A `master-config` ConfigMap is created in the `config-master` namespace containing application configuration in JSON format
2. **Automatic Reflection**: The Reflector operator automatically copies this ConfigMap to 4 child namespaces (`child-1`, `child-2`, `child-3`, `child-4`)
3. **Live Configuration Reading**: Go applications running in each child namespace continuously read the reflected configuration every 3 seconds
4. **Real-time Updates**: When you update the master ConfigMap, changes are automatically reflected to all child namespaces and picked up by the running applications

## Architecture

```
┌─────────────────┐    ┌──────────────────────────────────────┐
│  config-master  │    │           Reflector Operator         │
│                 │    │                                      │
│ master-config   │────┤  Watches for ConfigMap changes      │
│ ConfigMap       │    │  Reflects to target namespaces      │
└─────────────────┘    └──────────────────────────────────────┘
                                            │
                       ┌────────────────────┼────────────────────┐
                       │                    │                    │
                       ▼                    ▼                    ▼
               ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
               │   child-1    │    │   child-2    │    │   child-3    │
               │              │    │              │    │              │
               │ Go App Pod   │    │ Go App Pod   │    │ Go App Pod   │
               │ (reads every │    │ (reads every │    │ (reads every │
               │  3 seconds)  │    │  3 seconds)  │    │  3 seconds)  │
               └──────────────┘    └──────────────┘    └──────────────┘
                                            │
                                   ┌──────────────┐
                                   │   child-4    │
                                   │              │
                                   │ Go App Pod   │
                                   │ (reads every │
                                   │  3 seconds)  │
                                   └──────────────┘
```

## Components

### Go Application (`main.go`)
- Reads JSON configuration from a file every 3 seconds
- Accepts a `-config` flag to specify the configuration file path
- Prints current configuration to stdout on each read

### Configuration Structure
```json
{
  "port": 8080,
  "address": "0.0.0.0",
  "enabled": true
}
```

### Kubernetes Resources
- **1 Master Namespace**: `config-master` - Contains the source ConfigMap
- **4 Child Namespaces**: `child-1`, `child-2`, `child-3`, `child-4` - Receive reflected ConfigMaps
- **1 Master ConfigMap**: `master-config` - Source of truth for configuration
- **4 Application Pods**: One per child namespace, each mounting the reflected ConfigMap
- **Reflector Operator**: Handles automatic ConfigMap synchronization

## Prerequisites

- Kubernetes cluster (local or remote)
- [Skaffold](https://skaffold.dev/) installed
- Docker for building container images
- `kubectl` configured to access your cluster

## How to Run

1. **Clone and navigate to the project**:
   ```bash
   cd /path/to/ember-reflector-lab
   ```

2. **Start the demo**:
   ```bash
   skaffold dev
   ```

   This command will:
   - Build the Go application into a Docker container
   - Install the Emberstack Reflector operator
   - Create all namespaces
   - Deploy the master ConfigMap with reflection annotations
   - Start pods in each child namespace

3. **Watch the logs** to see the applications reading configuration:
   ```bash
   # Watch logs from all pods
   kubectl logs -f pod/config-reader-pod -n child-1
   kubectl logs -f pod/config-reader-pod -n child-2
   kubectl logs -f pod/config-reader-pod -n child-3
   kubectl logs -f pod/config-reader-pod -n child-4
   ```

4. **Test configuration reflection** by updating the master ConfigMap:
   ```bash
   kubectl patch configmap master-config -n config-master --type merge -p '{"data":{"config.json":"{\"port\":9090,\"address\":\"127.0.0.1\",\"enabled\":false}"}}'
   ```

   You should see the changes reflected in all child namespaces and picked up by the running applications within seconds.

## Key Features Demonstrated

- **Automatic ConfigMap Reflection**: Changes to the master ConfigMap are automatically synchronized
- **Namespace Isolation**: Each child namespace operates independently while sharing configuration
- **Real-time Configuration Updates**: Applications pick up configuration changes without restarts
- **Scalable Architecture**: Easy to add more child namespaces as needed

## Reflector Annotations

The master ConfigMap uses these annotations to enable reflection:
```yaml
annotations:
  reflector.v1.k8s.emberstack.com/reflection-allowed: "true"
  reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces: "child-1,child-2,child-3,child-4"
  reflector.v1.k8s.emberstack.com/reflection-auto-enabled: "true"
```

## Cleanup

To clean up all resources:
```bash
# Stop Skaffold (Ctrl+C) then run:
skaffold delete
```

This will remove all created namespaces, pods, ConfigMaps, and the Reflector operator.

## Troubleshooting

- **Pods not starting**: Check if the Docker image was built successfully
- **ConfigMap not reflecting**: Verify the Reflector operator is running in the `reflector-system` namespace
- **Application errors**: Check pod logs for configuration file read errors

## Learn More

- [Emberstack Reflector Documentation](https://github.com/emberstack/kubernetes-reflector)
- [Skaffold Documentation](https://skaffold.dev/docs/)
- [Kubernetes ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/)