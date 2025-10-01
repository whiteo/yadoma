// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package container

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/container"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) StopContainer(
	ctx context.Context,
	req *protos.StopContainerRequest,
) (*protos.StopContainerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}
	err := s.layer.StopContainer(ctx, req.GetId(), container.StopOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot stop container: %v", err)
	}

	return &protos.StopContainerResponse{Success: true}, nil
}
