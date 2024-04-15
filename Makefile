# This will build the docker image for ai-training-api
.PHONY: build-ai-training-api
build-ai-training-api:
	docker build -t grafana/metadata-service -f ./ai-training-api/Dockerfile .
