// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package volume provides the agent's service layer for Docker volume management.
// It implements gRPC-facing handlers that validate requests, delegate to the
// Docker client layer, map results to protobuf messages, and translate errors
// into gRPC status codes.
//
// Supported operations include creating volumes, listing and inspecting details,
// removing volumes, and pruning unused volumes. Calls respect the caller's
// context and deadlines; streaming endpoints are not used.
//
// The package does not spawn goroutines on behalf of the caller and relies on
// context cancellation for shutdown. It is intended for internal use by the
// agent's gRPC server layer.
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

// NewVolumeService constructs a new volume service backed by the provided Docker layer.
// It performs no I/O, starts no goroutines, and simply wires the service to the lower‑level client.
// The returned service can be registered on a gRPC server via Register.
// The provided layer should be non‑nil.
func NewVolumeService(layer *docker.Layer) *Service {
	return &Service{layer: layer}
}

// Register attaches the volume service handlers to the provided gRPC server.
// It exposes the service's RPC endpoints to clients and should be invoked
// during server setup, before the server starts serving requests.
// The call performs no network I/O and returns immediately.
func (s *Service) Register(rpc *grpc.Server) {
	protos.RegisterVolumeServiceServer(rpc, s)
}
