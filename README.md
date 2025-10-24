# Yadoma - Yet Another DOcker MAnager

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://golang.org)
![Java](https://img.shields.io/badge/Java-25-blue)
![Spring Boot](https://img.shields.io/badge/Spring_Boot-3.5.6-blue)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

[![codecov](https://codecov.io/gh/whiteo/yadoma/branch/main/graph/badge.svg)](https://codecov.io/gh/whiteo/yadoma)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=whiteo_yadoma2&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=whiteo_yadoma2)

[![Tests](https://github.com/whiteo/yadoma/actions/workflows/test.yml/badge.svg)](https://github.com/whiteo/yadoma/actions/workflows/test.yml)
[![Deploy Agent](https://github.com/whiteo/yadoma/actions/workflows/deploy-agent.yml/badge.svg)](https://github.com/whiteo/yadoma/actions/workflows/deploy-agent.yml)
[![Deploy WebApp](https://github.com/whiteo/yadoma/actions/workflows/deploy-cloudrun.yml/badge.svg)](https://github.com/whiteo/yadoma/actions/workflows/deploy-cloudrun.yml)

A modern Docker management platform with gRPC agent and web interface.

## Architecture

```
┌─────────────┐    REST/WS    ┌─────────────┐     gRPC      ┌─────────────┐
│  React UI   │◄─────────────►│   Spring    │◄─────────────►│   Agent     │
│   (Vite)    │               │    Boot     │               │   (gRPC)    │
└─────────────┘               └─────────────┘               └─────────────┘
                                     │                              │
                                     ▼                              ▼
                              ┌─────────────┐              ┌─────────────┐
                              │   MongoDB   │              │   Docker    │
                              └─────────────┘              └─────────────┘
```

**Components:**
- **Agent** - Lightweight Go service that manages Docker via gRPC
- **WebApp** - Spring Boot backend with REST API and WebSocket support
- **UI** - React frontend served from the webapp

## Features

**Agent:**
- Container lifecycle (create, start, stop, remove)
- Image management (pull, build, list, remove)
- Real-time logs and statistics
- Network and volume management
- System monitoring

**WebApp:**
- JWT authentication
- User and container management
- Real-time updates via WebSocket
- Role-based access control
- OpenAPI documentation

## Quick Start

### Local Development

**Agent:**
```bash
make build-agent
./build/yadoma-agent
```

**WebApp + UI:**
```bash
# Build UI and webapp together
make build-webapp

# Or run in development mode
cd ui && npm run dev          # UI on :3000
cd webapp && ./gradlew bootRun # Backend on :8080
```

### Docker

```bash
# Set version
export YADOMA_AGENT_VERSION=v0.1.0

# Start services
docker compose up -d
```

## Deployment

### Agent → GCP VM
Automatically deploys to GCP Compute Engine on push to `main`:
- Tests with coverage → Codecov
- Builds optimized binary
- Deploys via gcloud SSH

### WebApp + UI → Cloud Run
Automatically deploys to Google Cloud Run on push to `main`:
- Tests with coverage → Codecov
- SonarCloud analysis
- Builds proto → UI → backend
- Containerized deployment

## Configuration

### Agent
```bash
--agent-tcp-port=:50001    # gRPC listen address
--docker-socket=/var/run/docker.sock
```

### WebApp
Environment variables:
- `MONGO_URI` - MongoDB connection string
- `TOKEN_SECRET_KEY` - JWT signing key
- `TOKEN_EXPIRATION_TIME` - Token lifetime (ms)
- `GRPC_HOST` - Agent gRPC address
- `APP_LOG_LEVEL` - Logging level

## API

**Endpoints:**
- `/` - React UI
- `/swagger-ui/index.html` - API documentation
- `/actuator/health` - Health check
- `/api/*` - REST API
- `/ws/*` - WebSocket endpoints

**Authentication:**
```bash
POST /authenticate
{
  "username": "user",
  "password": "pass"
}
```

## Development

### Project Structure
```
yadoma/
├── agent/              # Go gRPC service
├── webapp/             # Spring Boot backend
├── ui/                 # React frontend
├── proto/              # Protocol buffers
└── .github/workflows/  # CI/CD pipelines
```

### Commands
```bash
# Build
make build-agent       # Build agent binary
make build-ui          # Build React UI
make build-webapp      # Build webapp with UI

# Test
make test-agent        # Test agent
make test-webapp       # Test webapp

# Clean
make clean             # Remove build artifacts
```

### Proto Generation
```bash
make tools             # Install protoc plugins
make generate          # Generate Go stubs
```

Webapp proto generation happens automatically via Gradle.

## Security

**Agent:**
- Non-root user
- Docker socket mount (requires careful host security)
- Internal network recommended

**WebApp:**
- JWT authentication
- BCrypt password hashing
- CORS configuration
- Role-based access control
- Container isolation per user

## License

MIT - See [LICENSE](LICENSE) for details.
