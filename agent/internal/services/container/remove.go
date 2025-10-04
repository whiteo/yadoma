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

// RemoveContainer deletes a Docker container by the request ID.
// It maps request flags to container.RemoveOptions (`Force`, `RemoveVolumes`) and
// invokes the Docker layer using the caller's context.
// On success, it returns `Success=true`; on failure, it returns gRPC status errors:
// `InvalidArgument` when the ID is empty, or `Internal` if the underlying removal fails.
// The call respects context deadlines and cancellation.
func (s *Service) RemoveContainer(
	ctx context.Context,
	req *protos.RemoveContainerRequest,
) (*protos.RemoveContainerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	opts := container.RemoveOptions{
		Force:         req.Force,
		RemoveVolumes: req.RemoveVolumes,
	}

	err := s.layer.RemoveContainer(ctx, req.GetId(), opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot remove container: %v", err)
	}

	return &protos.RemoveContainerResponse{Success: true}, nil
}
