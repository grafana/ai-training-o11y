# ai-training-observability

Layout is as follows:

Project root contains anything necessary to spin up a dev environment to test end-to-end
.config/ has config files for the dev environment
observability/ is the python exporter
ai-training-api/ contains the metadata (and maybe proxy) service and any necessary files for postgres (specifying database configuration, especially)
grafana-aitraining-app/ contains the grafana plugin