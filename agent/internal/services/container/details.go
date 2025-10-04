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

// GetContainerDetails retrieves detailed information about a Docker container by its ID.
// It validates the request (returns codes.InvalidArgument when the ID is empty) and delegates
// to the Docker layer using the provided context.
// On success, it maps the result to protos.GetContainerDetailsResponse, including status,
// created time, mounts, and network settings. On Docker layer failure, it returns a
// gRPC error with codes.Internal and additional context.
func (s *Service) GetContainerDetails(
	ctx context.Context,
	req *protos.GetContainerDetailsRequest,
) (*protos.GetContainerDetailsResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	details, err := s.layer.GetContainerDetails(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get container details: %v", err)
	}

	return &protos.GetContainerDetailsResponse{
		Id:       details.ID,
		Image:    details.Image,
		Name:     details.Name,
		Status:   extractStatus(details.State),
		Created:  details.Created,
		Mounts:   mapMounts(details.Mounts),
		Networks: mapNetworks(details.NetworkSettings.Networks),
	}, nil
}
