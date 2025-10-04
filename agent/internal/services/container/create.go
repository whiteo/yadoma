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

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateContainer creates a Docker container from the request.
// It validates that an image is provided, maps request fields into container, host, and networking configs,
// and delegates creation to the Docker layer with the given name and a default OCI platform.
// On success, it returns the new container ID.
// On failure, it returns a gRPC error (codes.InvalidArgument for missing image, codes.Internal for creation errors).
func (s *Service) CreateContainer(
	ctx context.Context,
	req *protos.CreateContainerRequest,
) (*protos.CreateContainerResponse, error) {
	if req.GetImage() == "" {
		return nil, status.Error(codes.InvalidArgument, "image is required")
	}

	config := mapConfig(req)
	hostConfig := mapHostConfig(req.GetHostConfig())
	networkingConfig := mapNetworking(req.GetNetworks())

	resp, err := s.layer.CreateContainer(ctx, config, hostConfig, networkingConfig, &ocispec.Platform{}, req.GetName())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create container: %v", err)
	}

	return &protos.CreateContainerResponse{Id: resp.ID}, nil
}
