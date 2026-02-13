# Fleet Event Stream

Production-ready Go microservice for processing GPS vehicle events at scale with Prometheus metrics and Kubernetes deployment.

## Quick Start

```bash
go build -o bin/fleet-event-stream ./cmd/api
./bin/fleet-event-stream
```

API runs on `:8080`, metrics on `:9090`.

## Endpoints

| Endpoint | Description |
|----------|-------------|
| `POST /api/v1/events` | Ingest vehicle event |
| `GET /api/v1/stats` | Event statistics |
| `GET /health` | Health check |
| `GET /ready` | Readiness probe |
| `GET :9090/metrics` | Prometheus metrics |

## Project Structure

```
cmd/api/          # Entry point
internal/
  handlers/       # HTTP handlers
  models/         # Data models & validation
  processor/      # Event processing logic
  metrics/        # Prometheus metrics
deployments/k8s/  # Kubernetes manifests
```

## Docker

```bash
docker build -t fleet-event-stream .
docker run -p 8080:8080 -p 9090:9090 fleet-event-stream
```

## Kubernetes

```bash
kubectl apply -f deployments/k8s/
```

## Configurations

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | API port |
| `METRICS_PORT` | `9090` | Metrics port |
