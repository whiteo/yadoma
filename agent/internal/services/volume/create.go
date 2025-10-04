// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package volume provides the agent's service layer for Docker volume management.
// It implements gRPC-facing handlers that validate requests, delegate to the
// Docker client layer, map results to protobuf messages, and translate errors
// into gRPC status codes.
//
// Supported operations include creating volumes, listing and inspecting details,
// removing volumes, and pruning unused volumes. Calls respect the caller's
// context and deadlines; streaming endpoints are not used.
//
// The package does not spawn goroutines on behalf of the caller and relies on
// context cancellation for shutdown. It is intended for internal use by the
// agent's gRPC server layer.
package volume

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/volume"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateVolume creates a Docker volume using fields from the request.
// It validates input (volume name is required), forwards the operation to the
// Docker client layer with the provided context, and maps the result to the
// protobuf response.
// On failure, it returns a gRPC error with an appropriate status code
// (codes.InvalidArgument for bad input, codes.Internal for client errors).
func (s *Service) CreateVolume(
	ctx context.Context,
	req *protos.CreateVolumeRequest,
) (*protos.CreateVolumeResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "volume name is required")
	}

	vol, err := s.layer.CreateVolume(ctx, volume.CreateOptions{
		Name:       req.Name,
		Driver:     req.Driver,
		DriverOpts: req.DriverOpts,
		Labels:     req.Labels,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create volume: %v", err)
	}

	return &protos.CreateVolumeResponse{
		Name:       vol.Name,
		Driver:     vol.Driver,
		Mountpoint: vol.Mountpoint,
		Labels:     vol.Labels,
		CreatedAt:  vol.CreatedAt,
		Scope:      vol.Scope,
	}, nil
}
