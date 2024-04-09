# ai-training-observability

Layout is as follows:

Project root contains anything necessary to spin up a dev environment to test end-to-end
.config/ has config files for the dev environment
observability/ is the python exporter
ai-training-metadata/ contains the metadata service and any necessary files for postgres (specifying database configuration, especially)
ai-training-observability/ contains the grafana plugin, or will once i move it over here

todo:

[ ] flesh out the design on the postgres tables regarding process + run tracking

[ ] dummy/basic workflows in the exporter to solidify what that flow looks like

[ ] initialize the metadata service

[ ] make a run work entirely between the metadata service and the exporter

[ ] make sure the plugin can connect to the metadata service, and especially can give back a login string to the user

[ ] dataviz stuff
