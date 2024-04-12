# Monitor ChromaDB

## Understand the Environment Variables

Firstly, here’s a quick overview of what each environment variable you'll be setting does:

1. **`CHROMA_OTEL_COLLECTION_ENDPOINT`**: This is the endpoint to which the metrics from ChromaDB will be sent. It uses the OpenTelemetry protocol to communicate with Grafana Cloud.

2. **`CHROMA_OTEL_SERVICE_NAME`**: The name for your ChromaDB service that will appear in Grafana to identify your metrics source.

3. **`CHROMA_OTEL_COLLECTION_HEADERS`**: These are the headers required for authentication with your Grafana Cloud instance. This particular one includes an authorization header with basic auth credentials.

4. **`CHROMA_OTEL_GRANULARITY`**: This defines the granularity of the telemetry data collected; in this case, all data points are collected.

## Configure Environment Variables

These variables need to be set in the environment where ChromaDB runs. How you do this depends on your deployment method. Here are general instructions for a few common scenarios:

### Direct Run

If you're running ChromaDB directly on a host (for development or testing), you can export these variables in your terminal:

```bash
export CHROMA_OTEL_COLLECTION_ENDPOINT=https://otlp-gateway[asdfasd].grafana.net/otlp
export CHROMA_OTEL_SERVICE_NAME=chromadb
export CHROMA_OTEL_COLLECTION_HEADERS=[asdfasdfsadf]
export CHROMA_OTEL_GRANULARITY=All
```

### Docker Container

For a Docker deployment, these variables can be added to your `docker run` command or your Docker Compose file:

```shell
docker run -e CHROMA_OTEL_COLLECTION_ENDPOINT=https://otlp-gateway[asdasdf]grafana.net/otlp \
  -e CHROMA_OTEL_SERVICE_NAME=chromadb \
  -e CHROMA_OTEL_COLLECTION_HEADERS=Authorization=[asdfasfsdf] \
  -e CHROMA_OTEL_GRANULARITY=All \
  your_chromadb_image
```

### Kubernetes

If you're deploying on Kubernetes, these variables should be added to your pod specifications under `env` in the deployment YAML:

```yaml
env:
  - name: CHROMA_OTEL_COLLECTION_ENDPOINT
    value: "https://[asdfsdf].grafana.net/otlp"
  - name: CHROMA_OTEL_SERVICE_NAME
    value: "chromadb"
  - name: CHROMA_OTEL_COLLECTION_HEADERS
    value: "Authorization="
  - name: CHROMA_OTEL_GRANULARITY
    value: "All"
```

## Install dashboard
Get access to pre-configured dashboard that work right away

