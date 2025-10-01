// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package network

import (
	"context"

	"github.com/whiteo/yadoma/internal/dockers"
	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"

	"google.golang.org/grpc"
)

type layerAPI interface {
	GetNetworks(ctx context.Context, opts network.ListOptions) ([]network.Summary, error)
	GetNetworkDetails(ctx context.Context, id string, opts network.InspectOptions) (network.Inspect, error)
	CreateNetwork(ctx context.Context, name string, opts network.CreateOptions) (network.CreateResponse, error)
	RemoveNetwork(ctx context.Context, id string) error
	ConnectNetwork(ctx context.Context, networkID, containerID string, config *network.EndpointSettings) error
	DisconnectNetwork(ctx context.Context, networkID, containerID string, force bool) error
	PruneNetworks(ctx context.Context, args filters.Args) (network.PruneReport, error)
}

type Service struct {
	protos.UnimplementedNetworkServiceServer
	layer layerAPI
}

func NewNetworkService(layer *docker.Layer) *Service {
	return &Service{layer: layer}
}

func (s *Service) Register(rpc *grpc.Server) {
	protos.RegisterNetworkServiceServer(rpc, s)
}
