# This will build the docker image for ai-training-api
.PHONY: build-ai-training-api
build-ai-training-api:
	docker build -t grafana/ai-training-api -f ./ai-training-api/Dockerfile .
