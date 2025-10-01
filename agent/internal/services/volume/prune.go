// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package volume

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/filters"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) PruneVolumes(
	ctx context.Context,
	req *protos.PruneVolumesRequest,
) (*protos.PruneVolumesResponse, error) {
	return s.pruneVolumes(ctx, req)
}

func (s *Service) pruneVolumes(
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
