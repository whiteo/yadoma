// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package system provides the agent's service layer for Docker system operations.
// It exposes gRPC-facing handlers that validate requests, delegate to the Docker
// client layer, map results to protobuf messages, and translate errors into gRPC
// status codes.
//
// Supported operations include retrieving daemon/system information and reporting
// aggregate disk usage across images, containers, volumes, and layer sizes.
// Calls respect the caller's context and deadlines; streaming endpoints are not used.
//
// The package does not spawn goroutines on behalf of the caller and relies on
// context cancellation for shutdown. It is intended for internal use by the
// agent's gRPC server layer.
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

// NewSystemService constructs a System service backed by the provided Docker layer.
// It sets the layer as the service's dependency to perform Docker system operations
// and returns a service ready to be registered via Register.
// The function starts no goroutines and performs no validation; callers should pass
// a non-nil layer and manage its lifecycle externally.
func NewSystemService(layer *docker.Layer) *Service {
	return &Service{layer: layer}
}

// Register attaches the System service to the provided gRPC server.
// It registers the service with the server so that RPCs defined in the
// protos.SystemService API become available to clients.
// The call starts no goroutines and transfers no ownership; the caller
// remains responsible for the server's lifecycle. Invoke this once per
// *grpc.Server instance; repeated registration will panic.
func (s *Service) Register(rpc *grpc.Server) {
	protos.RegisterSystemServiceServer(rpc, s)
}
