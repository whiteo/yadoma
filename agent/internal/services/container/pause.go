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

// PauseContainer pauses a Docker container identified by the provided ID using the service's container layer.
// It validates the request and returns a gRPC InvalidArgument error if the ID is missing.
// The call honors the incoming context; cancellation or timeout aborts the operation.
// On success, it returns a response with Success=true; on failure, it maps errors to a gRPC Internal status.
func (s *Service) PauseContainer(
	ctx context.Context,
	req *protos.PauseContainerRequest,
) (*protos.PauseContainerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	err := s.layer.PauseContainer(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot pause container: %v", err)
	}

	return &protos.PauseContainerResponse{Success: true}, nil
}
