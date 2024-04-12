#  Monitor TensorFlow Serving 

## Install Grafana Agent

The Grafana agent collects observability data and sends it to Grafana Cloud. Once the agent is deployed to your hosts, it collects and sends Prometheus-style metrics and log data using a pared-down Prometheus collector.

Run this command to install and run Grafana Agent as a grafana-agent.service systemd service

```shell
ARCH="amd64" GCLOUD_HOSTED_METRICS_URL="https://prometheus-prod-13-prod-us-east-0.grafana.net/api/prom/push" GCLOUD_HOSTED_METRICS_ID="" GCLOUD_SCRAPE_INTERVAL="60s" GCLOUD_HOSTED_LOGS_URL="https://logs-prod-006.grafana.net/loki/api/v1/push" GCLOUD_HOSTED_LOGS_ID="" GCLOUD_RW_API_KEY="" /bin/sh -c "$(curl -fsSL https://storage.googleapis.com/cloud-onboarding/agent/scripts/static/install-linux.sh)"
```

## Prepare your agent configuration file

### Logs
Below `logs.configs.scrape_configs`, insert the following lines according to your environment.

```
- job_name: integrations/tensorflow
  relabel_configs:
    - source_labels: ['__meta_docker_container_name']
      replacement: tensorflow
      target_label: name
    - source_labels: ['__meta_docker_container_name']
      replacement: integrations/tensorflow
      target_label: job
    - source_labels: ['__meta_docker_container_name']
      replacement: '<your-instance-name>'
      target_label: instance
  docker_sd_configs:
    - host: unix:///var/run/docker.sock
      refresh_interval: 5s
      filters:
        - name: name
          values: [tensorflow]
```

### Metrics
Below `metrics.configs.scrape_configs`, insert the following lines and change the URLs according to your environment:

```
- job_name: integrations/tensorflow
  metrics_path: "/monitoring/prometheus/metrics"
  relabel_configs:
    - replacement: '<your-instance-name>'
      target_label: instance
  static_configs:
    - targets: ['localhost:8501']
  metric_relabel_configs:
  - action: keep
    regex: :tensorflow:core:graph_build_calls|:tensorflow:core:graph_build_time_usecs|:tensorflow:core:graph_run_time_usecs|:tensorflow:core:graph_runs|:tensorflow:serving:batching_session:queuing_latency_count|:tensorflow:serving:batching_session:queuing_latency_sum|:tensorflow:serving:request_count|:tensorflow:serving:request_latency_count|:tensorflow:serving:request_latency_sum|:tensorflow:serving:runtime_latency_count|:tensorflow:serving:runtime_latency_sum
    source_labels:
      - __name__
```

## Restart the agent

Once you’ve made changes to your agent configuration file, run the following command to restart the agent.

After installation, the Agent config is stored in /etc/grafana-agent.yaml. Restart the agent for any changes to take effect:

```shell
sudo systemctl restart grafana-agent.service
```

## Install dashboard
Get access to pre-configured dashboard that work right away

```
{
  "annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": {
          "type": "grafana",
          "uid": "-- Grafana --"
        },
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
      }
    ]
  },
  "description": "Overview of a TensorFlow Serving instance.",
  "editable": true,
  "fiscalYearStartMonth": 0,
  "graphTooltip": 0,
  "id": 13,
  "links": [],
  "panels": [
    {
      "datasource": {
        "uid": "${prometheus_datasource}"
      },
      "description": "Rate of requests over time for the selected model. Grouped by statuses.",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "reqps"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 24,
        "x": 0,
        "y": 0
      },
      "id": 2,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "uid": "${prometheus_datasource}"
          },
          "expr": "rate(:tensorflow:serving:request_count{job=~\"$job\",instance=~\"$instance\",model_name=~\"$model_name\"}[$__rate_interval])",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "model_name=\"{{model_name}}\",status=\"{{status}}\"",
          "refId": "A"
        }
      ],
      "title": "Model request rate",
      "type": "timeseries"
    },
    {
      "datasource": {
        "uid": "${prometheus_datasource}"
      },
      "description": "Average request latency of predict requests for the selected model.",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "µs"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 8
      },
      "id": 3,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "uid": "${prometheus_datasource}"
          },
          "expr": "increase(:tensorflow:serving:request_latency_sum{job=~\"$job\",instance=~\"$instance\",model_name=~\"$model_name\"}[$__rate_interval])/increase(:tensorflow:serving:request_latency_count{job=~\"$job\",instance=~\"$instance\",model_name=~\"$model_name\"}[$__rate_interval])",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "model_name=\"{{model_name}}\"",
          "refId": "A"
        }
      ],
      "title": "Model predict request latency",
      "type": "timeseries"
    },
    {
      "datasource": {
        "uid": "${prometheus_datasource}"
      },
      "description": "Average runtime latency to fulfill a predict request for the selected model.",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "µs"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 8
      },
      "id": 4,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "uid": "${prometheus_datasource}"
          },
          "expr": "increase(:tensorflow:serving:runtime_latency_sum{job=~\"$job\",instance=~\"$instance\",model_name=~\"$model_name\"}[$__rate_interval])/increase(:tensorflow:serving:runtime_latency_count{job=~\"$job\",instance=~\"$instance\",model_name=~\"$model_name\"}[$__rate_interval])",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "model_name=\"{{model_name}}\"",
          "refId": "A"
        }
      ],
      "title": "Model predict runtime latency",
      "type": "timeseries"
    },
    {
      "collapsed": false,
      "datasource": {
        "type": "prometheus",
        "uid": "grafanacloud-prom"
      },
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 16
      },
      "id": 5,
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "grafanacloud-prom"
          },
          "refId": "A"
        }
      ],
      "title": "Serving overview",
      "type": "row"
    },
    {
      "datasource": {
        "uid": "${prometheus_datasource}"
      },
      "description": "Number of times TensorFlow Serving has created a new client graph.",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "calls"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 17
      },
      "id": 6,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": false
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "uid": "${prometheus_datasource}"
          },
          "expr": "increase(:tensorflow:core:graph_build_calls{job=~\"$job\",instance=~\"$instance\"}[$__rate_interval])",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "A"
        }
      ],
      "title": "Graph build calls",
      "type": "timeseries"
    },
    {
      "datasource": {
        "uid": "${prometheus_datasource}"
      },
      "description": "Number of graph executions.",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "runs"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 17
      },
      "id": 7,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": false
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "uid": "${prometheus_datasource}"
          },
          "expr": "increase(:tensorflow:core:graph_runs{job=~\"$job\",instance=~\"$instance\"}[$__rate_interval])",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "A"
        }
      ],
      "title": "Graph runs",
      "type": "timeseries"
    },
    {
      "datasource": {
        "uid": "${prometheus_datasource}"
      },
      "description": "Amount of time Tensorflow has spent creating new client graphs.",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "unit": "µs"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 25
      },
      "id": 8,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": false
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "uid": "${prometheus_datasource}"
          },
          "expr": "increase(:tensorflow:core:graph_build_time_usecs{job=~\"$job\",instance=~\"$instance\"}[$__rate_interval])/increase(:tensorflow:core:graph_build_calls{job=~\"$job\",instance=~\"$instance\"}[$__rate_interval])",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "A"
        }
      ],
      "title": "Graph build time",
      "type": "timeseries"
    },
    {
      "datasource": {
        "uid": "${prometheus_datasource}"
      },
      "description": "Amount of time spent executing graphs.",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "unit": "µs"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 25
      },
      "id": 9,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": false
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "uid": "${prometheus_datasource}"
          },
          "expr": "increase(:tensorflow:core:graph_run_time_usecs{job=~\"$job\",instance=~\"$instance\"}[$__rate_interval])/increase(:tensorflow:core:graph_runs{job=~\"$job\",instance=~\"$instance\"}[$__rate_interval])",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "A"
        }
      ],
      "title": "Graph run time",
      "type": "timeseries"
    },
    {
      "datasource": {
        "uid": "${prometheus_datasource}"
      },
      "description": "Current latency in the batching queue.",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "unit": "µs"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 33
      },
      "id": 10,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": false
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "uid": "${prometheus_datasource}"
          },
          "expr": "increase(:tensorflow:serving:batching_session:queuing_latency_sum{job=~\"$job\",instance=~\"$instance\"}[$__rate_interval])/increase(:tensorflow:serving:batching_session:queuing_latency_count{job=~\"$job\",instance=~\"$instance\"}[$__rate_interval])",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "A"
        }
      ],
      "title": "Batch queuing latency",
      "type": "timeseries"
    },
    {
      "datasource": {
        "uid": "${prometheus_datasource}"
      },
      "description": "Rate of batch queue throughput over time.",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "unit": "batches/s"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 33
      },
      "id": 11,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": false
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "uid": "${prometheus_datasource}"
          },
          "expr": "rate(:tensorflow:serving:batching_session:queuing_latency_count{job=~\"$job\",instance=~\"$instance\"}[$__rate_interval])",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "A"
        }
      ],
      "title": "Batch queue throughput",
      "type": "timeseries"
    },
    {
      "datasource": {
        "uid": "${loki_datasource}"
      },
      "description": "Logs from the TensorFlow Serving Docker container.",
      "gridPos": {
        "h": 8,
        "w": 24,
        "x": 0,
        "y": 41
      },
      "id": 12,
      "options": {
        "dedupStrategy": "none",
        "enableLogDetails": true,
        "prettifyLogMessage": false,
        "showCommonLabels": false,
        "showLabels": false,
        "showTime": true,
        "sortOrder": "Descending",
        "wrapLogMessage": false
      },
      "targets": [
        {
          "datasource": {
            "uid": "${loki_datasource}"
          },
          "editorMode": "code",
          "expr": "{name=\"tensorflow\",job=~\"$job\",instance=~\"$instance\"}",
          "legendFormat": "",
          "queryType": "range",
          "refId": "A"
        }
      ],
      "title": "Container logs",
      "type": "logs"
    }
  ],
  "refresh": "1m",
  "schemaVersion": 39,
  "tags": [
    "tensorflow-integration"
  ],
  "templating": {
    "list": [
      {
        "current": {
          "selected": false,
          "text": "grafanacloud-aio11y-prom",
          "value": "grafanacloud-prom"
        },
        "hide": 0,
        "includeAll": false,
        "label": "Data source",
        "multi": false,
        "name": "prometheus_datasource",
        "options": [],
        "query": "prometheus",
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "type": "datasource"
      },
      {
        "allValue": ".+",
        "current": {
          "selected": false,
          "text": "All",
          "value": "$__all"
        },
        "datasource": {
          "uid": "${prometheus_datasource}"
        },
        "definition": "",
        "hide": 0,
        "includeAll": true,
        "label": "Job",
        "multi": true,
        "name": "job",
        "options": [],
        "query": "label_values(:tensorflow:serving:request_count{}, job)",
        "refresh": 2,
        "regex": "",
        "skipUrlSync": false,
        "sort": 1,
        "tagValuesQuery": "",
        "tagsQuery": "",
        "type": "query",
        "useTags": false
      },
      {
        "allValue": ".+",
        "current": {
          "selected": false,
          "text": "All",
          "value": "$__all"
        },
        "datasource": {
          "uid": "${prometheus_datasource}"
        },
        "definition": "",
        "hide": 0,
        "includeAll": true,
        "label": "Instance",
        "multi": true,
        "name": "instance",
        "options": [],
        "query": "label_values(:tensorflow:serving:request_count{job=~\"$job\"}, instance)",
        "refresh": 2,
        "regex": "",
        "skipUrlSync": false,
        "sort": 1,
        "tagValuesQuery": "",
        "tagsQuery": "",
        "type": "query",
        "useTags": false
      },
      {
        "allValue": ".+",
        "current": {
          "selected": false,
          "text": "All",
          "value": "$__all"
        },
        "datasource": {
          "uid": "${prometheus_datasource}"
        },
        "definition": "",
        "hide": 0,
        "includeAll": true,
        "label": "Model name",
        "multi": false,
        "name": "model_name",
        "options": [],
        "query": "label_values(:tensorflow:serving:request_count{job=~\"$job\",instance=~\"$instance\"}, model_name)",
        "refresh": 2,
        "regex": "",
        "skipUrlSync": false,
        "sort": 1,
        "tagValuesQuery": "",
        "tagsQuery": "",
        "type": "query",
        "useTags": false
      },
      {
        "current": {
          "selected": false,
          "text": "grafanacloud-aio11y-alert-state-history",
          "value": "grafanacloud-alert-state-history"
        },
        "hide": 0,
        "includeAll": false,
        "label": "Loki datasource",
        "multi": false,
        "name": "loki_datasource",
        "options": [],
        "query": "loki",
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "type": "datasource"
      }
    ]
  },
  "time": {
    "from": "now-30m",
    "to": "now"
  },
  "timepicker": {
    "refresh_intervals": [
      "5s",
      "10s",
      "30s",
      "1m",
      "5m",
      "15m",
      "30m",
      "1h",
      "2h",
      "1d"
    ],
    "time_options": [
      "5m",
      "15m",
      "1h",
      "6h",
      "12h",
      "24h",
      "2d",
      "7d",
      "30d"
    ]
  },
  "timezone": "default",
  "title": "TensorFlow Serving overview",
  "uid": "tensorflow-overview",
  "version": 1,
  "weekStart": ""
}
```
=======
1. Install Grafana Agent```
<Same steps as we see in Integrations>, For hackathon maybe we put Linux instrcutions```

2. Update Agent Configuration

Add the following `scrape_config` to Grafana Agent Cofiguration

logs
    - job_name: integrations/tensorflow
      relabel_configs:
        - source_labels: ['__meta_docker_container_name']
          replacement: tensorflow
          target_label: name
        - source_labels: ['__meta_docker_container_name']
          replacement: integrations/tensorflow
          target_label: job
        - source_labels: ['__meta_docker_container_name']
          replacement: '<your-instance-name>'
          target_label: instance
      docker_sd_configs:
        - host: unix:///var/run/docker.sock
          refresh_interval: 5s
          filters:
            - name: name
              values: [tensorflow]
metrics
    - job_name: integrations/tensorflow
      metrics_path: "/monitoring/prometheus/metrics"
      relabel_configs:
        - replacement: '<your-instance-name>'
          target_label: instance
      static_configs:
        - targets: ['localhost:8501']
      metric_relabel_configs:
      - action: keep
        regex: :tensorflow:core:graph_build_calls|:tensorflow:core:graph_build_time_usecs|:tensorflow:core:graph_run_time_usecs|:tensorflow:core:graph_runs|:tensorflow:serving:batching_session:queuing_latency_count|:tensorflow:serving:batching_session:queuing_latency_sum|:tensorflow:serving:request_count|:tensorflow:serving:request_latency_count|:tensorflow:serving:request_latency_sum|:tensorflow:serving:runtime_latency_count|:tensorflow:serving:runtime_latency_sum
        source_labels:
            - __ name __
           
