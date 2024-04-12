#  Monitor PyTorch TorchServe 

## Install Grafana Agent

The Grafana agent collects observability data and sends it to Grafana Cloud. Once the agent is deployed to your hosts, it collects and sends Prometheus-style metrics and log data using a pared-down Prometheus collector.

Run this command to install and run Grafana Agent as a grafana-agent.service systemd service

```shell
ARCH="amd64" GCLOUD_HOSTED_METRICS_URL="https://prometheus-prod-13-prod-us-east-0.grafana.net/api/prom/push" GCLOUD_HOSTED_METRICS_ID="1444735" GCLOUD_SCRAPE_INTERVAL="60s" GCLOUD_HOSTED_LOGS_URL="https://logs-prod-006.grafana.net/loki/api/v1/push" GCLOUD_HOSTED_LOGS_ID="821357" GCLOUD_RW_API_KEY="" /bin/sh -c "$(curl -fsSL https://storage.googleapis.com/cloud-onboarding/agent/scripts/static/install-linux.sh)"
```

## Prepare your agent configuration file

### Metrics

Below `metrics.configs.scrape_configs`, insert the following lines and change the URLs according to your environment:

```
- job_name: 'torchserve'
  static_configs:
  - targets: ['localhost:8082'] #TorchServe metrics endpoint
```

## Restart the agent

Once you’ve made changes to your agent configuration file, run the following command to restart the agent.

After installation, the Agent config is stored in /etc/grafana-agent.yaml. Restart the agent for any changes to take effect:

```shell
sudo systemctl restart grafana-agent.service
```

## Install dashboard
Get access to pre-configured dashboard that work right away

