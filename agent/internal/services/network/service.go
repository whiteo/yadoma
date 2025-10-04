// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package network provides the agent's service layer for Docker network management.
// It exposes gRPC-facing handlers that validate requests, delegate to the Docker
// layer, map results to protobuf messages, and translate errors into gRPC status
// codes.
//
// Supported operations include creating networks, listing and inspecting details,
// connecting and disconnecting containers, removing networks, and pruning unused
// networks. Calls respect the caller's context and deadlines; streaming endpoints
// are not used.
//
// The package does not spawn goroutines on behalf of the caller and relies on
// context cancellation for shutdown. It is intended for internal use by the
// agent's gRPC server layer.
package network

import (
	"context"

	docker "github.com/whiteo/yadoma/internal/dockers"
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

// NewNetworkService creates a network service backed by the given Docker layer.
// It wires the layer into the gRPC-facing service and returns an instance ready
// to be registered via Register. The service does not spawn background goroutines
// and relies on caller-provided context for cancellation.
func NewNetworkService(layer *docker.Layer) *Service {
	return &Service{layer: layer}
}

// Register attaches the NetworkService implementation to the provided gRPC server.
// It exposes all RPC endpoints generated from the protobuf definition by calling
// `protos.RegisterNetworkServiceServer`.
// Call this once during server setup before starting to serve requests.
// The method is non-blocking, starts no background goroutines, and relies on the
// server's lifecycle for shutdown. Passing a nil server will cause a panic.
func (s *Service) Register(rpc *grpc.Server) {
	protos.RegisterNetworkServiceServer(rpc, s)
}
