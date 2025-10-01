// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package container

import (
	"context"
	"io"

	docker "github.com/whiteo/yadoma/internal/dockers"
	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"google.golang.org/grpc"
)

type layerAPI interface {
	GetContainers(ctx context.Context, opts container.ListOptions) ([]container.Summary, error)
	GetContainerDetails(ctx context.Context, id string) (container.InspectResponse, error)
	GetContainerLogs(ctx context.Context, id string, opts container.LogsOptions) (io.ReadCloser, error)
	GetContainerStats(ctx context.Context, id string, stream bool) (container.StatsResponseReader, error)
	CreateContainer(ctx context.Context,
		config *container.Config,
		hostConfig *container.HostConfig,
		networkingConfig *network.NetworkingConfig,
		platform *ocispec.Platform,
		containerName string,
	) (container.CreateResponse, error)
	RemoveContainer(ctx context.Context, id string, opts container.RemoveOptions) error
	StartContainer(ctx context.Context, id string, opts container.StartOptions) error
	StopContainer(ctx context.Context, id string, opts container.StopOptions) error
	RestartContainer(ctx context.Context, id string, opts container.StopOptions) error
	PauseContainer(ctx context.Context, id string) error
	UnpauseContainer(ctx context.Context, id string) error
	KillContainer(ctx context.Context, id, signal string) error
	RenameContainer(ctx context.Context, id, name string) error
}

type Service struct {
	protos.UnimplementedContainerServiceServer
	layer layerAPI
}

func NewContainerService(layer *docker.Layer) *Service {
	return &Service{layer: layer}
}

func (s *Service) Register(rpc *grpc.Server) {
	protos.RegisterContainerServiceServer(rpc, s)
}
