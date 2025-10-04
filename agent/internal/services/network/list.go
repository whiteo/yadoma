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

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/network"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetNetworks lists Docker networks using default network.ListOptions.
// It honors the incoming context for cancellation and deadlines and maps results
// to protos.GetNetworksResponse with basic network fields.
// On failure, it returns a gRPC status error with code codes.Internal.
func (s *Service) GetNetworks(
	ctx context.Context,
	_ *protos.GetNetworksRequest,
) (*protos.GetNetworksResponse, error) {
	list, err := s.layer.GetNetworks(ctx, network.ListOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot list networks: %v", err)
	}

	resp := &protos.GetNetworksResponse{
		Networks: make([]*protos.GetNetworkResponse, 0, len(list)),
	}

	for _, n := range list {
		resp.Networks = append(resp.Networks, &protos.GetNetworkResponse{
			Id:     n.ID,
			Name:   n.Name,
			Driver: n.Driver,
			Scope:  n.Scope,
		})
	}

	return resp, nil
}
