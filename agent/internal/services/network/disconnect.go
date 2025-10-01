// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package network

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) DisconnectNetwork(
	ctx context.Context,
	req *protos.DisconnectNetworkRequest,
) (*protos.DisconnectNetworkResponse, error) {
	if req.GetNetworkId() == "" {
		return nil, status.Error(codes.InvalidArgument, "network ID is required")
	}

	if req.GetContainerId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	err := s.layer.DisconnectNetwork(ctx, req.GetNetworkId(), req.GetContainerId(), req.GetForce())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot disconnect network: %v", err)
	}

	return &protos.DisconnectNetworkResponse{}, nil
}
