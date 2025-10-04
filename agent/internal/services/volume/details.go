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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetVolumeDetails returns details for a Docker volume by its ID.
// It validates the request, respects the caller\'s context for cancellation,
// delegates to the underlying client layer, and maps the result to a protobuf response.
// On failure, it returns a gRPC status error (e.g., InvalidArgument for an empty ID, Internal on client errors).
func (s *Service) GetVolumeDetails(
	ctx context.Context,
	req *protos.GetVolumeDetailsRequest,
) (*protos.GetVolumeDetailsResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "volume ID is required")
	}

	details, err := s.layer.GetVolumeDetails(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get volume details: %v", err)
	}

	return &protos.GetVolumeDetailsResponse{
		Name:       details.Name,
		Driver:     details.Driver,
		Mountpoint: details.Mountpoint,
		Labels:     details.Labels,
		CreatedAt:  details.CreatedAt,
		Scope:      details.Scope,
	}, nil
}
