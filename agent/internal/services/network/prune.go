// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package network

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/filters"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) PruneNetworks(
	ctx context.Context,
	req *protos.PruneNetworksRequest,
) (*protos.PruneNetworksResponse, error) {
	r, err := s.layer.PruneNetworks(ctx, filters.NewArgs(filters.Arg("All", req.All)))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot prune network: %v", err)
	}

	return &protos.PruneNetworksResponse{
		NetworksDeleted: r.NetworksDeleted,
	}, nil
}
