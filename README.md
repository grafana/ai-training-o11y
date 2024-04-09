# ai-training-observability

Layout is as follows:

Project root contains anything necessary to spin up a dev environment to test end-to-end
.config/ has config files for the dev environment
observability/ is the python exporter
ai-training-metadata/ contains the metadata service and any necessary files for postgres (specifying database configuration, especially)
ai-training-observability/ contains the grafana plugin, or will once i move it over here
