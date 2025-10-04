// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package container provides service-layer operations for managing Docker containers.
// It implements gRPC-facing logic that validates requests, invokes the Docker layer,
// maps results to protobuf messages, and returns errors as gRPC status codes.
// Supported operations cover the container lifecycle and inspection, including create,
// list, inspect, logs and stats streaming, start/stop/restart, kill, pause/unpause,
// rename, and remove. Calls respect the caller's context; streaming endpoints propagate
// cancellation and require the caller to consume and close streams. The package is
// internal to the agent and intended to be used by higher-level gRPC servers.
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

// NewContainerService constructs a new Service backed by the provided Docker layer.
// It initializes the service dependency used to perform container operations and
// returns an instance ready to be registered on a gRPC server via Register.
// Callers should provide a non-nil layer to avoid runtime failures.
func NewContainerService(layer *docker.Layer) *Service {
	return &Service{layer: layer}
}

// Register attaches this service implementation to the provided gRPC server.
// It wires up the generated ContainerServiceServer handlers so RPCs are exposed
// to clients. Call during server initialization before starting to serve.
func (s *Service) Register(rpc *grpc.Server) {
	protos.RegisterContainerServiceServer(rpc, s)
}
