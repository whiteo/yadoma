// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package system

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetDiskUsage(
	ctx context.Context,
	_ *protos.GetDiskUsageRequest,
) (*protos.GetDiskUsageResponse, error) {
	usage, err := s.layer.GetDiskUsage(ctx, types.DiskUsageOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get disk usage info: %v", err)
	}

	images := mapDiskUsageImage(usage.Images)
	containers := mapDiskUsageContainer(usage.Containers)
	volumes := mapDiskUsageVolume(usage.Volumes)

	return &protos.GetDiskUsageResponse{
		LayersSize: usage.LayersSize,
		Images:     images,
		Containers: containers,
		Volumes:    volumes,
	}, nil
}
