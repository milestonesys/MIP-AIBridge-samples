#!/bin/bash
set -eo pipefail

# Setup image configuration
export TARGETARCH=amd64

# Build the docker image
COMPOSE_BAKE=true DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 docker compose build
