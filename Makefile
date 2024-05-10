# This will build the docker image for ai-training-api
build-ai-training-api:
	docker build --target development -t grafana/ai-training-api -f ./ai-training-api/Dockerfile .

build-ai-training-api-prod:
	docker build -t grafana/ai-training-api -f ./ai-training-api/Dockerfile .

## Calls "mage" in the ai-training-app directory to build the app
.phony: build-aitraining-app
build-aitraining-app:
	cd grafana-aitraining-app && mage -v
	cd grafana-aitraining-app && yarn install
	cd grafana-aitraining-app && yarn build

# Brings up grafana, associated databases, and the ai-training-api
.PHONY: docker
docker: build-ai-training-api build-aitraining-app
	docker compose up

# Self explanatory
.PHONY: docker-down
docker-down:
	docker compose down

# Builds the exporter for whatever your processor architecture is and linux for docker
.PHONY: jupyter-exporter
jupyter-exporter:
	@HOST_ARCH=$$(uname -m); \
	if [ "$$HOST_ARCH" = "x86_64" ]; then \
		TARGET_ARCH=amd64; \
	elif [ "$$HOST_ARCH" = "aarch64" ] || [ "$$HOST_ARCH" = "arm64" ]; then \
		TARGET_ARCH=arm64; \
	else \
		echo "Unsupported architecture: $$HOST_ARCH"; \
		exit 1; \
	fi; \
	cd o11y && TARGET_OS=linux TARGET_ARCH=$$TARGET_ARCH hatch build -t wheel

# Builds every conceivable wheel
# We may need to add more of these eventually because some HPC clusters run outdated architectures
.PHONY: exporter-wheels
exporter-wheels:
	@for os in linux mac windows; do \
		for arch in amd64 arm64; do \
			echo "Building for OS: $$os, Architecture: $$arch"; \
			cd o11y && TARGET_OS=$$os TARGET_ARCH=$$arch hatch build -t wheel; \
			cd ..; \
		done; \
	done

# Brings up jupyter so you can run smoke.ipynb and any related workflows
.PHONY: jupyter
jupyter: jupyter-exporter
	docker compose -f docker-compose-jupyter.yml up

# This won't work on macs
.PHONY: jupytorch
jupytorch: exporter-build
	docker compose -f docker-compose-pytorch.yml up

.PHONY: clean
clean:
	cd o11y && rm -rf __pycache__ .cache .ipython .ipynb_checkpoints .jupyter .local .npm dist
	cd o11y/src/go-plugin && rm -rf dist
# This only cleans the o11y build folders, not the other two modules
