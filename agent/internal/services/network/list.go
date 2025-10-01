// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package network

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/network"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetNetworks(
	ctx context.Context,
	_ *protos.GetNetworksRequest,
) (*protos.GetNetworksResponse, error) {
	list, err := s.layer.GetNetworks(ctx, network.ListOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot list networks: %v", err)
	}

	resp := &protos.GetNetworksResponse{
		Networks: make([]*protos.GetNetworkResponse, 0, len(list)),
	}

	for _, n := range list {
		resp.Networks = append(resp.Networks, &protos.GetNetworkResponse{
			Id:     n.ID,
			Name:   n.Name,
			Driver: n.Driver,
			Scope:  n.Scope,
		})
	}

	return resp, nil
}
