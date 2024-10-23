# syntax=docker/dockerfile:1.4
# go-sqlite3 requires CGO_ENABLED=1 so we need to use a build image

ARG GO_VERSION=1.22

# The ai-training-api folder does not have any git files
# so we need to copy the .git folder from the root of the repository
FROM golang:${GO_VERSION}-bullseye as git_build_meta
WORKDIR /git_meta
COPY --link .git ./.git
RUN git describe --tags --long --always --abbrev=10 > .git_version && \
  git rev-parse HEAD^{commit} > .git_commit && \
  git rev-parse --abbrev-ref HEAD > .git_branch

FROM golang:${GO_VERSION}-bullseye as prep
# Avoid 'detected dubious ownership in repository' errors caused by differing
# UID ownership when building.
RUN git config --global --add safe.directory /go/src

ENV GOCACHE=/opt/go/gocache \
  GOMODCACHE=/opt/go/gomodcache

WORKDIR /go/src/ai-training-api

COPY --link ai-training-api/go.* ./
COPY --from=git_build_meta /git_meta/.git_* ./
COPY --link ai-training-api ./

FROM prep AS development

# Air fixed at 1.52.3 because 1.60.0 would force upgrading go to 1.23 and this is easier
RUN --mount=type=cache,id=go-cache-ai-training-api,target=/opt/go go install github.com/air-verse/air@v1.52.3

FROM prep as build

RUN --mount=type=cache,id=go-cache-ai-training-api,target=/opt/go \
  BUILD_VERSION=$(cat .git_version) \
  BUILD_COMMIT=$(cat .git_commit) \
  BUILD_BRANCH=$(cat .git_branch) \
  make exe

FROM debian:bullseye-slim as production

RUN apt-get update && apt-get -y install ca-certificates \
  && rm -rf /var/lib/apt/lists/*

COPY --link --from=build \
  go/src/ai-training-api/dist/ai-training-api \
  /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/ai-training-api"]
