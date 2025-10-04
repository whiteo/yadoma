// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package container provides service-layer operations for managing Docker containers.
// It implements gRPC-facing logic that validates requests, invokes the Docker layer,
// maps results to protobuf messages, and returns errors as gRPC status codes.
// Supported operations cover the container lifecycle and inspection, including create,
// list, inspect, logs and stats streaming, start/stop/restart, kill, pause/unpause,
// rename, and remove. Calls respect the caller's context; streaming endpoints propagate
// cancellation and require the caller to consume and close streams. The package is
// internal to the agent and intended to be used by higher-level gRPC servers.
package container

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/container"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RestartContainer restarts a Docker container by its ID via the Docker layer.
// It validates that the request contains a non-empty container ID
// and honors the caller's context for cancellation and deadlines.
// Errors are mapped to gRPC status codes: InvalidArgument for an empty ID and Internal if the restart fails.
// On success, it returns a protos.RestartContainerResponse with Success set to true.
func (s *Service) RestartContainer(
	ctx context.Context,
	req *protos.RestartContainerRequest,
) (*protos.RestartContainerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	if err := s.layer.RestartContainer(ctx, req.GetId(), container.StopOptions{}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to restart container: %v", err)
	}

	return &protos.RestartContainerResponse{Success: true}, nil
}
