// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package container

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/container"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) RemoveContainer(
	ctx context.Context,
	req *protos.RemoveContainerRequest,
) (*protos.RemoveContainerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	opts := container.RemoveOptions{
		Force:         req.Force,
		RemoveVolumes: req.RemoveVolumes,
	}

	err := s.layer.RemoveContainer(ctx, req.GetId(), opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot remove container: %v", err)
	}

	return &protos.RemoveContainerResponse{Success: true}, nil
}
