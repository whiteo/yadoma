// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package container

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/container"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) RestartContainer(
	ctx context.Context,
	req *protos.RestartContainerRequest,
) (*protos.RestartContainerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	if err := s.layer.RestartContainer(ctx, req.GetId(), container.StopOptions{}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to restart container: %v", err)
	}

	return &protos.RestartContainerResponse{Success: true}, nil
}
