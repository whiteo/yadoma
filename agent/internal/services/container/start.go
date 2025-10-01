// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package container

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/container"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) StartContainer(
	ctx context.Context,
	req *protos.StartContainerRequest,
) (*protos.StartContainerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	err := s.layer.StartContainer(ctx, req.GetId(), container.StartOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot start container: %v", err)
	}

	return &protos.StartContainerResponse{Success: true}, nil
}
