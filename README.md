# Yadoma - Yet Another DOcker MAnager

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

[![Docker Image Version](https://img.shields.io/docker/v/whiteo/yadoma?label=docker&sort=semver)](https://hub.docker.com/r/whiteo/yadoma)

[![Go Report Card](https://goreportcard.com/badge/github.com/whiteo/yadoma)](https://goreportcard.com/report/github.com/whiteo/yadoma)
[![codecov](https://codecov.io/gh/whiteo/yadoma/branch/master/graph/badge.svg)](https://codecov.io/gh/whiteo/yadoma)

[![Tests](https://github.com/whiteo/yadoma/actions/workflows/test.yml/badge.svg)](https://github.com/whiteo/yadoma/actions/workflows/test.yml)


[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=whiteo_yadoma&metric=security_rating&token=2994c76c1451733892c105b6694b54397a691daa)](https://sonarcloud.io/summary/new_code?id=whiteo_yadoma)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=whiteo_yadoma&metric=sqale_rating&token=2994c76c1451733892c105b6694b54397a691daa)](https://sonarcloud.io/summary/new_code?id=whiteo_yadoma)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=whiteo_yadoma&metric=reliability_rating&token=2994c76c1451733892c105b6694b54397a691daa)](https://sonarcloud.io/summary/new_code?id=whiteo_yadoma)

## Overview

Yadoma is a lightweight agent for managing Docker containers via gRPC API.

## Features

- **Container Management**: Create, start, stop, remove, logs, stats
- **Image Operations**: Pull, build, list, remove, prune  
- **Network Control**: Create, connect, disconnect, remove networks
- **Volume Management**: Create, list, remove, prune volumes
- **System Information**: Disk usage, system info monitoring

## Architecture

- **gRPC Server**: High-performance API with Protocol Buffers
- **Docker Integration**: Direct Docker Engine API communication
- **Modular Services**: Separate services for each Docker resource type
- **Structured Logging**: JSON logging with zerolog

## Quick Start

```bash
# Build the agent
make build

# Run the gRPC server
./bin/yadoma-agent

# Default port: 50051
```

## Development

```bash
# Generate protobuf files
make generate

# Run tests
make test

# Build for production
make build
```
