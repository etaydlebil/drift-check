# drift-check

> CLI tool that detects configuration drift between running Kubernetes workloads and their source Helm charts.

---

## Installation

```bash
go install github.com/yourusername/drift-check@latest
```

Or download a pre-built binary from the [Releases](https://github.com/yourusername/drift-check/releases) page.

---

## Usage

Point `drift-check` at a Helm chart and a running namespace to surface any differences between the deployed workloads and the chart's rendered manifests.

```bash
# Check for drift between a deployed release and its local chart
drift-check --release my-app --namespace production --chart ./charts/my-app

# Use a remote chart from a Helm repository
drift-check --release my-app --namespace production --repo https://charts.example.com --chart my-app --version 1.4.2
```

### Flags

| Flag | Description |
|-------------|----------------------------------------------|
| `--release` | Name of the deployed Helm release |
| `--namespace` | Kubernetes namespace of the release |
| `--chart` | Path or name of the source Helm chart |
| `--repo` | Helm chart repository URL (optional) |
| `--version` | Chart version to compare against (optional) |
| `--output` | Output format: `text`, `json`, or `yaml` |

### Example Output

```
[DRIFT DETECTED] deployment/my-app
  - spec.replicas: expected 3, got 1
  - spec.template.spec.containers[0].image: expected app:1.4.2, got app:1.3.0

[OK] service/my-app
[OK] configmap/my-app-config
```

---

## Requirements

- Go 1.21+
- `kubectl` configured with access to your cluster
- Helm 3

---

## License

This project is licensed under the [MIT License](LICENSE).