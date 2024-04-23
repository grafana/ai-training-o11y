# This will build the docker image for ai-training-api
.PHONY: build-ai-training-api
build-ai-training-api:
	$(MAKE) exe -C ai-training-api
	docker build -t grafana/ai-training-api -f ./ai-training-api/Dockerfile .

## Calls "mage" in the ai-training-app directory to build the app
.phony: build-aitraining-app
build-aitraining-app:
	cd grafana-aitraining-app && mage -v
	cd grafana-aitraining-app && yarn install
	cd grafana-aitraining-app && yarn build	

.PHONY: docker
docker: build-ai-training-api build-aitraining-app
	docker compose up

.PHONY: docker-down
docker
	docker compose down