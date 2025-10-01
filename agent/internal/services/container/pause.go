// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package container

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) PauseContainer(
	ctx context.Context,
	req *protos.PauseContainerRequest,
) (*protos.PauseContainerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	err := s.layer.PauseContainer(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot pause container: %v", err)
	}

	return &protos.PauseContainerResponse{Success: true}, nil
}
