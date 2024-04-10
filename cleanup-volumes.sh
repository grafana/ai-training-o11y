#!/bin/bash

# Stop and remove containers
docker compose down

# Remove volumes (postgres data)
# This is desirable if changing how the postgres database is initialized, or if it is large.
docker volume rm ai-training-observability_postgres_data