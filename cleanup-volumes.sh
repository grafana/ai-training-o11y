#!/bin/bash

# Stop and remove containers
docker compose down

# Remove volumes (postgres)
docker volume rm ai-training-observability_postgres_data