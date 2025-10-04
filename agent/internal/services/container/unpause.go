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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnpauseContainer unpauses a Docker container identified by the request ID.
// It validates that the request contains a non\-empty container ID and delegates
// the operation to the underlying Docker layer using the provided context.
// On success, it returns a response with `Success=true`.
// On failure, it returns a gRPC error: `InvalidArgument` if the ID is empty,
// or `Internal` if the Docker layer returns an error.
func (s *Service) UnpauseContainer(
	ctx context.Context,
	req *protos.UnpauseContainerRequest,
) (*protos.UnpauseContainerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	err := s.layer.UnpauseContainer(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot unpause container: %v", err)
	}

	return &protos.UnpauseContainerResponse{Success: true}, nil
}
