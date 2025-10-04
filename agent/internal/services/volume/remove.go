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

// RemoveVolume removes a Docker volume by its ID via the service's client layer.
// It validates the request ID, respects the incoming context for cancellation/deadlines,
// and delegates the operation to the underlying layer with the provided force flag.
// On success, it returns a response with `Success` set to true.
// On failure, it returns gRPC errors: `InvalidArgument` for a missing ID and `Internal` for client failures.
func (s *Service) RemoveVolume(
	ctx context.Context,
	req *protos.RemoveVolumeRequest,
) (*protos.RemoveVolumeResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "volume ID is required")
	}

	if err := s.layer.RemoveVolume(ctx, req.Id, req.Force); err != nil {
		return nil, status.Errorf(codes.Internal, "cannot remove volume: %v", err)
	}

	return &protos.RemoveVolumeResponse{
		Success: true,
	}, nil
}
