// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package container

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) KillContainer(
	ctx context.Context,
	req *protos.KillContainerRequest,
) (*protos.KillContainerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	err := s.layer.KillContainer(ctx, req.GetId(), req.GetSignal())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot kill container: %v", err)
	}

	return &protos.KillContainerResponse{Success: true}, nil
}
