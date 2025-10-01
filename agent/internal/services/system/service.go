// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package system

import (
	"context"

	docker "github.com/whiteo/yadoma/internal/dockers"
	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/system"

	"google.golang.org/grpc"
)

type layerAPI interface {
	GetSystemInfo(ctx context.Context) (system.Info, error)
	GetDiskUsage(ctx context.Context, opts types.DiskUsageOptions) (types.DiskUsage, error)
}

type Service struct {
	protos.UnimplementedSystemServiceServer
	layer layerAPI
}

func NewSystemService(layer *docker.Layer) *Service {
	return &Service{layer: layer}
}

func (s *Service) Register(rpc *grpc.Server) {
	protos.RegisterSystemServiceServer(rpc, s)
}
