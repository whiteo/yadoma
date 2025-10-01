// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package volume

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) RemoveVolume(
	ctx context.Context,
	req *protos.RemoveVolumeRequest,
) (*protos.RemoveVolumeResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "volume ID is required")
	}

	if err := s.layer.RemoveVolume(ctx, req.Id, req.Force); err != nil {
		return nil, status.Errorf(codes.Internal, "cannot remove volume: %v", err)
	}

	return &protos.RemoveVolumeResponse{
		Success: true,
	}, nil
}
