// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package volume

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/volume"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) CreateVolume(
	ctx context.Context,
	req *protos.CreateVolumeRequest,
) (*protos.CreateVolumeResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "volume name is required")
	}

	vol, err := s.layer.CreateVolume(ctx, volume.CreateOptions{
		Name:       req.Name,
		Driver:     req.Driver,
		DriverOpts: req.DriverOpts,
		Labels:     req.Labels,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create volume: %v", err)
	}

	return &protos.CreateVolumeResponse{
		Name:       vol.Name,
		Driver:     vol.Driver,
		Mountpoint: vol.Mountpoint,
		Labels:     vol.Labels,
		CreatedAt:  vol.CreatedAt,
		Scope:      vol.Scope,
	}, nil
}
