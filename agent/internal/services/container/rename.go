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

// RenameContainer renames an existing Docker container to the provided `name`.
// It validates that the request contains a non-empty container ID and delegates
// the operation to the Docker layer, propagating the caller's context for cancellation.
// On success, it returns a `RenameContainerResponse` with `Success` set to true.
// On failure, it returns gRPC errors: `InvalidArgument` for a missing ID and `Internal` for Docker-layer failures.
func (s *Service) RenameContainer(
	ctx context.Context,
	req *protos.RenameContainerRequest,
) (*protos.RenameContainerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	err := s.layer.RenameContainer(ctx, req.GetId(), req.GetName())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot rename container: %v", err)
	}

	return &protos.RenameContainerResponse{Success: true}, nil
}
