// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package volume

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/volume"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetVolumes(ctx context.Context, _ *protos.GetVolumesRequest) (*protos.GetVolumesResponse, error) {
	vols, err := s.layer.GetVolumes(ctx, volume.ListOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot list volumes: %v", err)
	}

	resp := &protos.GetVolumesResponse{
		Volumes:  make([]*protos.Volume, 0, len(vols.Volumes)),
		Warnings: vols.Warnings,
	}

	for _, v := range vols.Volumes {
		resp.Volumes = append(resp.Volumes, &protos.Volume{
			Name:       v.Name,
			Driver:     v.Driver,
			Mountpoint: v.Mountpoint,
			Labels:     v.Labels,
		})
	}

	return resp, nil
}
