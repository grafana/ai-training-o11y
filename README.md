# AI Training Observability

Layout is as follows:

- Project root contains anything necessary to spin up a dev environment to test end-to-end
- `.config/` has config files for the dev environment
- `o11y/` is the python exporter
- `ai-training-api/` contains the metadata (and maybe proxy) service and any necessary files for mysql (specifying database configuration, especially)
- `grafana-aitraining-app/` contains the grafana plugin

## Development
Requires:
- Python (3.8 or later)
- Hatch (best installed via "pipx install hatch" if you have pipx)
- Node (20 or later)
- Go (1.22 or later)
- Docker
- Yarn
- Make
- Mage

Builds dev environment with "make docker"

Once it's up, you can background it or open another terminal and use "make jupyter" to open a jupyter server. It will have a link to jupyter in your terminal.

Grafana will be hosted at localhost:3000 with the plugin.
