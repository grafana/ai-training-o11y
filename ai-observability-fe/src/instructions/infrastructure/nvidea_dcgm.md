# Monitor Nvidea DCGM

## Install Grafana Agent in Kubernetes

The Grafana agent collects observability data and sends it to Grafana Cloud. Once the agent is deployed to your hosts, it collects and sends Prometheus-style metrics and log data using a pared-down Prometheus collector.

Run this command to install and run Grafana Agent in Kubernetes

```
helm repo add grafana https://grafana.github.io/helm-charts &&
  helm repo update &&
  helm upgrade --install --atomic --timeout 300s grafana-k8s-monitoring grafana/k8s-monitoring \
    --namespace "default" --create-namespace --values - <<EOF
cluster:
  name: my-cluster
externalServices:
  prometheus:
    host: https://prometheus[url].grafana.net
    basicAuth:
      username: 
      password: 
  loki:
    host: https://logs[url].grafana.net
    basicAuth:
      username: 
      password: 
  tempo:
    host: https://tempo-prod-04-prod-us-east-0.grafana.net:443
    basicAuth:
      username: 
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
        url: https://prometheus-prod-13-prod-us-east-0.grafana.net/api/prom
kube-state-metrics:
  enabled: true
prometheus-node-exporter:
  enabled: true
prometheus-operator-crds:
  enabled: true
grafana-agent: {}
grafana-agent-logs: {}
EOF
```

## Install DCGM Exporter   

The NVIDIA DCGM Exporter fetches metrics from GPUs and exposes them for collection. It's crucial for monitoring the performance and health of your GPUs within Kubernetes.

```
  helm repo add gpu-helm-charts https://nvidia.github.io/dcgm-exporter/helm-charts

  helm repo update

  helm install --generate-name gpu-helm-charts/dcgm-exporter
```

## Prepare your agent configuration file

Below `metrics.configs.scrape_configs`, insert the following lines:  
    
```yml  
    scrape_configs:
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
        target_label: kubernetes_node  
```


## Install Dashboard
