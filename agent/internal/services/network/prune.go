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

	"github.com/docker/docker/api/types/filters"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PruneNetworks prunes Docker networks according to the request.
// It respects the incoming context for cancellation and deadlines and delegates
// to the Docker layer, passing an "All" filter derived from req.All.
// On success, it returns identifiers of deleted networks; on failure, a gRPC
// error with codes.Internal is returned.
func (s *Service) PruneNetworks(
	ctx context.Context,
	req *protos.PruneNetworksRequest,
) (*protos.PruneNetworksResponse, error) {
	r, err := s.layer.PruneNetworks(ctx, filters.NewArgs(filters.Arg("All", req.All)))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot prune network: %v", err)
	}

	return &protos.PruneNetworksResponse{
		NetworksDeleted: r.NetworksDeleted,
	}, nil
}
