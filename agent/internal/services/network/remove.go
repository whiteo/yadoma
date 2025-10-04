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

// RemoveNetwork removes a Docker network identified by the provided ID.
// It validates that the request contains a non-empty network ID, then delegates
// the deletion to the underlying Docker layer using the caller's context.
// On success, it returns an empty RemoveNetworkResponse.
// Errors are translated to gRPC status codes: InvalidArgument if the ID is missing,
// and Internal if the removal operation fails.
func (s *Service) RemoveNetwork(
	ctx context.Context,
	req *protos.RemoveNetworkRequest,
) (*protos.RemoveNetworkResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "network ID is required")
	}

	err := s.layer.RemoveNetwork(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot remove network: %v", err)
	}

	return &protos.RemoveNetworkResponse{}, nil
}
