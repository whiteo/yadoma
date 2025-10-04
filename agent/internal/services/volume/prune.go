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

	"github.com/docker/docker/api/types/filters"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PruneVolumes prunes unused Docker volumes based on the incoming request.
// It delegates to the Docker client layer with the provided context, passing a
// filter built from req.All (true: prune all unused volumes, false: dangling only).
// On success, it returns the names of deleted volumes and the total reclaimed bytes.
// On failure, it returns a gRPC Internal status wrapping the underlying error.
func (s *Service) PruneVolumes(
	ctx context.Context,
	req *protos.PruneVolumesRequest,
) (*protos.PruneVolumesResponse, error) {
	r, err := s.layer.PruneVolumes(ctx, filters.NewArgs(filters.Arg("All", req.All)))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot prune volume: %v", err)
	}

	return &protos.PruneVolumesResponse{
		VolumesDeleted: r.VolumesDeleted,
		SpaceReclaimed: r.SpaceReclaimed,
	}, nil
}
