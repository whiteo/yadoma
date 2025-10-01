// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package container

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) UnpauseContainer(
	ctx context.Context,
	req *protos.UnpauseContainerRequest,
) (*protos.UnpauseContainerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	err := s.layer.UnpauseContainer(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot unpause container: %v", err)
	}

	return &protos.UnpauseContainerResponse{Success: true}, nil
}
