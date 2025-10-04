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

// StartContainer starts a Docker container by ID.
// It validates the request, then calls s.layer.StartContainer with default container.StartOptions.
// Returns gRPC errors: codes.InvalidArgument if ID is missing, codes.Internal on start failure.
// On success, returns protos.StartContainerResponse with Success=true.
// The call honors ctx; cancellation and deadlines are propagated to the Docker client.
func (s *Service) StartContainer(
	ctx context.Context,
	req *protos.StartContainerRequest,
) (*protos.StartContainerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	err := s.layer.StartContainer(ctx, req.GetId(), container.StartOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot start container: %v", err)
	}

	return &protos.StartContainerResponse{Success: true}, nil
}
