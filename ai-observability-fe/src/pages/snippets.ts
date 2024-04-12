// THIS MARKS THE INFRA SNIPPETS SECTION
export const dcgmSnippet = () => {
    return `helm repo add grafana https://grafana.github.io/helm-charts &&
    helm repo update &&
    helm upgrade --install --atomic --timeout 300s grafana-k8s-monitoring grafana/k8s-monitoring \
      --namespace "default" --create-namespace --values - <<EOF
  cluster:
    name: my-cluster
  externalServices:
    prometheus:
      host: https://prometheus-[url].grafana.net
      basicAuth:
        username: ""
        password: 
      host: https://logs-[url].grafana.net
      basicAuth:
        username: ""
        password:
    tempo:
      host: https://tempo-[url]grafana.net
      basicAuth:
        username: ""
        password: 
  metrics:
    enabled: true
    cost:
      enabled: true
    node-exporter:
      enabled: true
  logs:
    enabled: true
    pod_logs:
      enabled: true
    cluster_events:
      enabled: true
  traces:
    enabled: true
  opencost:
    enabled: true
    opencost:
      exporter:
        defaultClusterId: my-cluster
      prometheus:
        external:
          url: https://prometheus-[url].grafana.net/api/prom
  kube-state-metrics:
    enabled: true
  prometheus-node-exporter:
    enabled: true
  prometheus-operator-crds:
    enabled: true
  grafana-agent: {}
  grafana-agent-logs: {}
  EOF`
}

export const installDCGMExporter = () => {
  return  `  helm repo add gpu-helm-charts https://nvidia.github.io/dcgm-exporter/helm-charts
  helm repo update
  helm install --generate-name gpu-helm-charts/dcgm-exporter`
}

export const prepareAgentConfig = () => {
  return `scrape_configs:
  job_name: gpu-metrics
 scrape_interval: 1s
 metrics_path: /metrics
 scheme: http
 kubernetes_sd_configs:
  role: endpoints
     namespaces:
     names:
      gpu-operator
 relabel_configs:
  source_labels: [__meta_kubernetes_pod_node_name]
     action: replace
     target_label: kubernetes_node `
}
//  THIS ENDS THE INFRA SNIPPETS SECTION

// THIS MARKS THE LLM SNIPPETS SECTION
// THIS ENDS LLM SNIPPETS SECTION


// THIS MARKS THE VECTOR DB SECTION
export const chromaDirectRun = () => {
  return`
  export CHROMA_OTEL_COLLECTION_ENDPOINT=https://otlp-gateway-prod-us-east-0.grafana.net/otlp
  export CHROMA_OTEL_SERVICE_NAME=chromadb
  export CHROMA_OTEL_COLLECTION_HEADERS=Authorization=Basic%20ODY1MjM0OmdsY19leUp2SWpvaU5qVXlPVGt5SWl3aWJpSTZJbk4wWVdOckxUZzJOVEl6TkMxdmRHeHdMWGR5YVhSbExYZGhiR3hoYUdraUxDSnJJam9pVUZCR01EY3lTbnB2TmxKd1JqWnVaak15VmpVMFVGUTBJaXdpYlNJNmV5SnlJam9pY0hKdlpDMTFjeTFsWVhOMExUQWlmWDA9
  export CHROMA_OTEL_GRANULARITY=All`

}

export const chromaDockerContainer = () => {
  return `
docker run -e CHROMA_OTEL_COLLECTION_ENDPOINT=https://otlp-gateway-prod-us-east-0.grafana.net/otlp \
  -e CHROMA_OTEL_SERVICE_NAME=chromadb \
  -e CHROMA_OTEL_COLLECTION_HEADERS=Authorization=Basic%20ODY1MjM0OmdsY19leUp2SWpvaU5qVXlPVGt5SWl3aWJpSTZJbk4wWVdOckxUZzJOVEl6TkMxdmRHeHdMWGR5YVhSbExYZGhiR3hoYUdraUxDSnJJam9pVUZCR01EY3lTbnB2TmxKd1JqWnVaak15VmpVMFVGUTBJaXdpYlNJNmV5SnlJam9pY0hKdlpDMTFjeTFsWVhOMExUQWlmWDA9 \
  -e CHROMA_OTEL_GRANULARITY=All \
  your_chromadb_image`
}

export const chromaKubernetes = () => {
return`env:
- name: CHROMA_OTEL_COLLECTION_ENDPOINT
  value: "https://otlp-gateway-prod-us-east-0.grafana.net/otlp"
- name: CHROMA_OTEL_SERVICE_NAME
  value: "chromadb"
- name: CHROMA_OTEL_COLLECTION_HEADERS
  value: "Authorization=Basic%20ODY1MjM0OmdsY19leUp2SWpvaU5qVXlPVGt5SWl3aWJpSTZJbk4wWVdOckxUZzJOVEl6TkMxdmRHeHdMWGR5YVhSbExYZGhiR3hoYUdraUxDSnJJam9pVUZCR01EY3lTbnB2TmxKd1JqWnVaak15VmpVMFVGUTBJaXdpYlNJNmV5SnlJam9pY0hKdlpDMTFjeTFsWVhOMExUQWlmWDA9"
- name: CHROMA_OTEL_GRANULARITY
  value: "All"`
}

export const pineconeDownloadAlloy = () => {
  return`ARCH="amd64" GCLOUD_HOSTED_METRICS_URL="https://prometheus[asdf]-prod-us-east-0.grafana.net/api/prom/push" GCLOUD_HOSTED_METRICS_ID="" GCLOUD_SCRAPE_INTERVAL="60s" GCLOUD_HOSTED_LOGS_URL="https://[asdfsad].grafana.net/loki/api/v1/push" GCLOUD_HOSTED_LOGS_ID="" GCLOUD_RW_API_KEY="" /bin/sh -c "$(curl -fsSL https://storage.googleapis.com/cloud-onboarding/agent/scripts/static/install-linux.sh)"`
}

export const pineconeMetrics = () => {
  return`- job_name: pinecone
  authorization:
    credentials: 
  scheme: https
  static_configs:
    - targets: ['metrics.YOUR_ENVIRONMENT.pinecone.io/metrics']`
}
// THIS ENDS THE VECTOR DB SECTION

// THIS BEGINS ML FRAMEWORKS SECTION


export const grafanaAlloyPytorchServe = () => {
  return `ARCH="amd64" GCLOUD_HOSTED_METRICS_URL="https://prometheus-prod-13-prod-us-east-0.grafana.net/api/prom/push" GCLOUD_HOSTED_METRICS_ID="1444735" GCLOUD_SCRAPE_INTERVAL="60s" GCLOUD_HOSTED_LOGS_URL="https://logs-prod-006.grafana.net/loki/api/v1/push" GCLOUD_HOSTED_LOGS_ID="821357" GCLOUD_RW_API_KEY="" /bin/sh -c "$(curl -fsSL https://storage.googleapis.com/cloud-onboarding/agent/scripts/static/install-linux.sh)"`
}

export const mlPytorchScrapeConfig = () => {
  return `- job_name: 'torchserve'
  static_configs:
  - targets: ['localhost:8082'] #TorchServe metrics endpoint`
}