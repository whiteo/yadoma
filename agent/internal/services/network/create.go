// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package network

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) CreateNetwork(
	ctx context.Context,
	req *protos.CreateNetworkRequest,
) (*protos.CreateNetworkResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "network name is required")
	}

	r, err := s.layer.CreateNetwork(ctx, req.Name, mapCreateOptions(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create network: %v", err)
	}

	return &protos.CreateNetworkResponse{Id: r.ID}, nil
}
