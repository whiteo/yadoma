// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package volume

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetVolumeDetails(
	ctx context.Context,
	req *protos.GetVolumeDetailsRequest,
) (*protos.GetVolumeDetailsResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "volume ID is required")
	}

	details, err := s.layer.GetVolumeDetails(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get volume details: %v", err)
	}

	return &protos.GetVolumeDetailsResponse{
		Name:       details.Name,
		Driver:     details.Driver,
		Mountpoint: details.Mountpoint,
		Labels:     details.Labels,
		CreatedAt:  details.CreatedAt,
		Scope:      details.Scope,
	}, nil
}
