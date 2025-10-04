# Yadoma - Yet Another DOcker MAnager

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

[![Docker Image Version](https://img.shields.io/docker/v/whiteo/yadoma?label=docker&sort=semver)](https://hub.docker.com/r/whiteo/yadoma)

[![Go Report Card](https://goreportcard.com/badge/github.com/whiteo/yadoma)](https://goreportcard.com/report/github.com/whiteo/yadoma)
[![codecov](https://codecov.io/gh/whiteo/yadoma/branch/master/graph/badge.svg)](https://codecov.io/gh/whiteo/yadoma)

[![Tests](https://github.com/whiteo/yadoma/actions/workflows/test.yml/badge.svg)](https://github.com/whiteo/yadoma/actions/workflows/test.yml)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=whiteo_yadoma2&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=whiteo_yadoma2)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=whiteo_yadoma2&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=whiteo_yadoma2)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=whiteo_yadoma2&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=whiteo_yadoma2)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=whiteo_yadoma2&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=whiteo_yadoma2)

## Overview

Yadoma is a lightweight agent for managing Docker containers via gRPC API.

## Features

- Container Management: create, start, stop, remove, logs, stats
- Image Operations: pull, build, list, remove, prune
- Network Control: create, connect, disconnect, remove networks
- Volume Management: create, list, remove, prune volumes
- System Information: disk usage, system info monitoring

## Quick Start

- Build the agent binary
  - `make build`
- Run the gRPC server locally
  - `./build/yadoma-agent --agent-tcp-port=:50001`
- Default gRPC port: 50001

### Docker Compose

- Use explicit image tags (no latest). The compose file respects env var `YADOMA_AGENT_VERSION`.
  - On Windows cmd.exe:
    - `set YADOMA_AGENT_VERSION=v0.1.0`
    - `docker compose up -d`
- To expose only locally, bind to loopback (see commented example in docker-compose.yml).

## Configuration: flags and environment

- Flags (agent/cmd/main.go):
  - `--agent-tcp-port`: string, default ":50001". gRPC listen address.
  - `--dockers-socket`: string, default "/var/run/docker.sock". Path to Docker Engine socket inside container/host mount.
- Environment:
  - `YADOMA_AGENT_VERSION`: image tag used by docker-compose.yml (e.g., v0.1.0).
  - `DOCKER_GID`: build arg passed to Dockerfile to align docker group GID inside the container.

## Image versioning (no latest)

- Make targets now build/push images with explicit tags:
  - `IMAGE_NAME`: defaults to `whiteo/yadoma-agent`.
  - `VERSION`: resolved from `YADOMA_AGENT_VERSION`, or `git describe`, or `v0.1.0` as fallback.
  - docker build: `make docker-build-agent` → builds `IMAGE_NAME:VERSION`.
  - docker push: `make docker-push-agent` → pushes `IMAGE_NAME:VERSION`.
- Compose uses `whiteo/yadoma-agent:${YADOMA_AGENT_VERSION:-v0.1.0}`.

## Tooling versions (protoc, plugins)

- Pinned versions in Makefile for reproducible codegen:
  - `protoc-gen-go`: `$(PROTOC_GEN_GO_VERSION)` (default v1.34.2)
  - `protoc-gen-go-grpc`: `$(PROTOC_GEN_GO_GRPC_VERSION)` (default v1.5.1)
- Install tools:
  - `make tools`
- Generate protobuf stubs:
  - `make generate`

## Security: mounting /var/run/docker.sock

Mounting the host Docker socket gives the container near-root control of the host. Treat it as highly privileged.

- Network segmentation:
  - Place the agent in an internal Docker network (see `agent_internal` in docker-compose.yml) and avoid exposing it publicly.
  - Prefer binding gRPC to 127.0.0.1 or a private interface only.
- Access control:
  - Limit who can reach the gRPC port (firewalls/security groups, reverse proxy mTLS if exposed).
  - Use per-environment credentials/ACLs on the entrypoint that talks to the agent (if any).
- Host hardening:
  - Run the container as non-root (Dockerfile already uses an unprivileged user in the docker group).
  - Keep Docker Engine updated; restrict who can access the host’s docker group.
  - Consider a dedicated host/VM segment for the agent if strict isolation is required.

## Development

- Generate protobuf files:
  - `make generate`
- Run tests:
  - `make test-agent`
- Build for production:
  - `make build`

## Logging

- Structured JSON logging via zerolog to stdout (see `agent/pkg/loggers/logger.go`).
- Timestamps are RFC3339; duration fields use milliseconds.
