// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package container

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) RenameContainer(
	ctx context.Context,
	req *protos.RenameContainerRequest,
) (*protos.RenameContainerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	err := s.layer.RenameContainer(ctx, req.GetId(), req.GetName())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot rename container: %v", err)
	}

	return &protos.RenameContainerResponse{Success: true}, nil
}
