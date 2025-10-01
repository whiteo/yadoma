// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package network

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) RemoveNetwork(
	ctx context.Context,
	req *protos.RemoveNetworkRequest,
) (*protos.RemoveNetworkResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "network ID is required")
	}

	err := s.layer.RemoveNetwork(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot remove network: %v", err)
	}

	return &protos.RemoveNetworkResponse{}, nil
}
