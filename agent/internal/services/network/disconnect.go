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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DisconnectNetwork disconnects a container from a Docker network.
// It validates that both network ID and container ID are provided,
// then delegates the operation to the underlying Docker layer with
// the provided force flag.
// The incoming context is propagated without modification.
// On success, returns an empty DisconnectNetworkResponse.
// On failure, returns a gRPC error: codes.InvalidArgument for missing
// inputs or codes.Internal if the disconnect operation fails.
func (s *Service) DisconnectNetwork(
	ctx context.Context,
	req *protos.DisconnectNetworkRequest,
) (*protos.DisconnectNetworkResponse, error) {
	if req.GetNetworkId() == "" {
		return nil, status.Error(codes.InvalidArgument, "network ID is required")
	}

	if req.GetContainerId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	err := s.layer.DisconnectNetwork(ctx, req.GetNetworkId(), req.GetContainerId(), req.GetForce())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot disconnect network: %v", err)
	}

	return &protos.DisconnectNetworkResponse{}, nil
}
