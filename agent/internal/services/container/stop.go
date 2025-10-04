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

// StopContainer stops a Docker container identified by req.Id through the Docker layer.
// It validates the request and returns gRPC InvalidArgument if the container ID is missing.
// The call honors the provided context and will abort if ctx is canceled or times out.
// On Docker-layer failure, it returns a gRPC Internal error with details.
// On success, it returns a StopContainerResponse with Success set to true.
func (s *Service) StopContainer(
	ctx context.Context,
	req *protos.StopContainerRequest,
) (*protos.StopContainerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}
	err := s.layer.StopContainer(ctx, req.GetId(), container.StopOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot stop container: %v", err)
	}

	return &protos.StopContainerResponse{Success: true}, nil
}
