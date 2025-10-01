// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package volume

import (
	"context"

	docker "github.com/whiteo/yadoma/internal/dockers"
	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"

	"google.golang.org/grpc"
)

type layerAPI interface {
	GetVolumes(ctx context.Context, opts volume.ListOptions) (volume.ListResponse, error)
	GetVolumeDetails(ctx context.Context, volumeID string) (volume.Volume, error)
	CreateVolume(ctx context.Context, opts volume.CreateOptions) (volume.Volume, error)
	RemoveVolume(ctx context.Context, volumeID string, force bool) error
	PruneVolumes(ctx context.Context, pruneFilters filters.Args) (volume.PruneReport, error)
}

type Service struct {
	protos.UnimplementedVolumeServiceServer
	layer layerAPI
}

func NewVolumeService(layer *docker.Layer) *Service {
	return &Service{layer: layer}
}

func (s *Service) Register(rpc *grpc.Server) {
	protos.RegisterVolumeServiceServer(rpc, s)
}
