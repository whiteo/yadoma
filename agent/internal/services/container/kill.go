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

// KillContainer sends the requested signal to the target container.
// It validates that the container ID is provided, propagates the incoming context
// for cancellation/timeouts, and delegates to the Docker layer with the container ID and signal.
// On success, it returns a response with Success=true.
// On error, it returns a gRPC status: codes.InvalidArgument if the ID is empty,
// or codes.Internal if the kill operation fails.
func (s *Service) KillContainer(
	ctx context.Context,
	req *protos.KillContainerRequest,
) (*protos.KillContainerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	err := s.layer.KillContainer(ctx, req.GetId(), req.GetSignal())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot kill container: %v", err)
	}

	return &protos.KillContainerResponse{Success: true}, nil
}
