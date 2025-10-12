# Yadoma - Yet Another DOcker MAnager

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://golang.org)
![java](https://img.shields.io/badge/Java-25-blue)
![spring](https://img.shields.io/badge/Spring_Boot-3.5.6-blue)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

[![Go Report Card](https://goreportcard.com/badge/github.com/whiteo/yadoma)](https://goreportcard.com/report/github.com/whiteo/yadoma)
[![codecov](https://codecov.io/gh/whiteo/yadoma/branch/master/graph/badge.svg)](https://codecov.io/gh/whiteo/yadoma)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=whiteo_yadoma2&metric=coverage)](https://sonarcloud.io/summary/new_code?id=whiteo_yadoma2)

[![Tests](https://github.com/whiteo/yadoma/actions/workflows/test.yml/badge.svg)](https://github.com/whiteo/yadoma/actions/workflows/test.yml)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=whiteo_yadoma2&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=whiteo_yadoma2)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=whiteo_yadoma2&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=whiteo_yadoma2)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=whiteo_yadoma2&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=whiteo_yadoma2)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=whiteo_yadoma2&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=whiteo_yadoma2)

## Overview

Yadoma is a comprehensive Docker management solution consisting of two main components:

1. **Agent** - A lightweight gRPC server for direct Docker daemon interaction
2. **WebApp** - A Spring Boot web application providing REST API and WebSocket interfaces for container management

## Architecture

```
┌─────────────────┐    gRPC     ┌─────────────────┐    Docker API    ┌─────────────────┐
│   Web Client   │◄───────────►│   Spring Boot   │◄────────────────►│  Yadoma Agent   │
│                 │  REST/WS    │     WebApp      │                  │   (gRPC Server) │
└─────────────────┘             └─────────────────┘                  └─────────────────┘
                                         │                                      │
                                         ▼                                      ▼
                                ┌─────────────────┐                   ┌─────────────────┐
                                │    MongoDB      │                   │  Docker Daemon  │
                                │   (Database)    │                   │   (Containers)  │
                                └─────────────────┘                   └─────────────────┘
```

## Features

### Agent Features
- Container Management: create, start, stop, remove, logs, stats
- Image Operations: pull, build, list, remove, prune
- Network Control: create, connect, disconnect, remove networks
- Volume Management: create, list, remove, prune volumes
- System Information: disk usage, system info monitoring

### WebApp Features
- **User Authentication**: JWT-based authentication and authorization
- **REST API**: Complete RESTful interface for container management
- **Real-time Updates**: WebSocket support for live container logs and statistics
- **User Management**: Multi-user support with role-based access control
- **Container Access Control**: Users can only access their assigned containers
- **Admin Interface**: Administrative functions for user and container management
- **Database Persistence**: MongoDB integration for user and container metadata storage

## Quick Start

### Agent
- Build the agent binary
  - `make build`
- Run the gRPC server locally
  - `./build/yadoma-agent --agent-tcp-port=:50001`
- Default gRPC port: 50001

### WebApp
- Navigate to webapp directory
  - `cd webapp`
- Run with Gradle
  - `./gradlew bootRun`
- Default web server port: 8080
- API documentation available at: `/swagger-ui.html`

### Docker Compose

- Use explicit image tags (no latest). The compose file respects env var `YADOMA_AGENT_VERSION`.
  - On Windows cmd.exe:
    - `set YADOMA_AGENT_VERSION=v0.1.0`
    - `docker compose up -d`
- To expose only locally, bind to loopback (see commented example in docker-compose.yml).

## Configuration

### Agent Configuration: flags and environment

- Flags (agent/cmd/main.go):
  - `--agent-tcp-port`: string, default ":50001". gRPC listen address.
  - `--dockers-socket`: string, default "/var/run/docker.sock". Path to Docker Engine socket inside container/host mount.
- Environment:
  - `YADOMA_AGENT_VERSION`: image tag used by docker-compose.yml (e.g., v0.1.0).
  - `DOCKER_GID`: build arg passed to Dockerfile to align docker group GID inside the container.

### WebApp Configuration

Configuration is managed through Spring Boot application properties:

- **Database**: MongoDB connection settings
- **Security**: JWT token configuration, CORS settings
- **gRPC Client**: Connection settings for communicating with the agent
- **WebSocket**: Real-time communication settings

Key configuration files:
- `webapp/src/main/resources/application.yml` - Main configuration
- `webapp/src/main/resources/application-dev.yml` - Development profile

## API Documentation

### REST Endpoints

- **Authentication**:
  - `POST /api/auth/login` - User login
  - `POST /api/auth/register` - User registration
  - `GET /api/auth/me` - Get current user info
  - `POST /api/auth/validate` - Validate JWT token

- **Container Management**:
  - `GET /api/containers` - List user's containers
  - `GET /api/containers/{id}` - Get container details
  - `POST /api/containers` - Create new container
  - `DELETE /api/containers/{id}` - Remove container
  - `POST /api/containers/{id}/start` - Start container
  - `POST /api/containers/{id}/stop` - Stop container
  - `POST /api/containers/{id}/restart` - Restart container

- **User Management** (Admin only):
  - `GET /api/users` - List all users
  - `POST /api/users` - Create new user
  - `PUT /api/users/{id}` - Update user

### WebSocket Endpoints

- `/ws/containers/{id}/logs?token={jwt}` - Real-time container logs
- `/ws/containers/{id}/stats?token={jwt}` - Real-time container statistics

## Testing

### Agent Tests
- Run agent tests:
  - `make test-agent`

### WebApp Tests
- Run webapp tests:
  - `cd webapp && ./gradlew test`
- Test coverage: **99.5%** (222 of 224 tests passing)
- Coverage report: `webapp/build/reports/jacoco/test/html/index.html`

### Test Coverage by Component:
- **Controllers**: 100%
- **Services**: 45% (business logic)
- **Security**: 83%
- **Exception Handling**: 94%
- **Data Transfer Objects**: 100%
- **Configuration**: 100%

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

## Security

### General Security: mounting /var/run/docker.sock

Mounting the host Docker socket gives the container near-root control of the host. Treat it as highly privileged.

- Network segmentation:
  - Place the agent in an internal Docker network (see `agent_internal` in docker-compose.yml) and avoid exposing it publicly.
  - Prefer binding gRPC to 127.0.0.1 or a private interface only.
- Access control:
  - Limit who can reach the gRPC port (firewalls/security groups, reverse proxy mTLS if exposed).
  - Use per-environment credentials/ACLs on the entrypoint that talks to the agent (if any).
- Host hardening:
  - Run the container as non-root (Dockerfile already uses an unprivileged user in the docker group).
  - Keep Docker Engine updated; restrict who can access the host's docker group.
  - Consider a dedicated host/VM segment for the agent if strict isolation is required.

### WebApp Security Features

- **JWT Authentication**: Secure token-based authentication
- **Role-based Access Control**: USER and ADMIN roles with different permissions
- **Container Isolation**: Users can only access containers assigned to them
- **Input Validation**: Comprehensive input validation and sanitization
- **CORS Configuration**: Configurable cross-origin resource sharing
- **Password Security**: Secure password hashing and validation
- **Session Management**: Secure session handling with JWT tokens

## Development

### Agent Development
- Generate protobuf files:
  - `make generate`
- Run tests:
  - `make test-agent`
- Build for production:
  - `make build`

### WebApp Development
- **Prerequisites**: Java 21, MongoDB
- **IDE Setup**: Project uses Lombok annotations
- **Development Profile**: Use `application-dev.yml` for local development
- **Database**: Embedded MongoDB for testing, external MongoDB for production
- **Hot Reload**: Use `./gradlew bootRun` for development with auto-restart

### Project Structure
```
webapp/
├── src/main/java/dev/whiteo/yadoma/
│   ├── config/          # Spring configuration classes
│   ├── controller/      # REST controllers
│   ├── domain/          # Domain entities
│   ├── dto/             # Data transfer objects
│   ├── exception/       # Exception handling
│   ├── mapper/          # Object mappers
│   ├── repository/      # Data repositories
│   ├── security/        # Security components
│   ├── service/         # Business logic services
│   ├── util/            # Utility classes
│   └── websocket/       # WebSocket endpoints
└── src/test/            # Test classes (mirror structure)
```

## Logging

### Agent Logging
- Structured JSON logging via zerolog to stdout (see `agent/pkg/loggers/logger.go`).
- Timestamps are RFC3339; duration fields use milliseconds.

### WebApp Logging
- Spring Boot default logging with configurable levels
- Structured logging for API requests and responses
- Error tracking and exception logging
- WebSocket connection and message logging
