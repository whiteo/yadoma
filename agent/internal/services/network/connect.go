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

// ConnectNetwork connects a container to a Docker network using settings from the request.
// It validates that both network ID and container ID are provided, then delegates to the Docker layer
// with endpoint settings mapped via mapEndpointSettings. The call honors the incoming context.
// On success, it returns an empty ConnectNetworkResponse. On failure, it returns gRPC errors:
// codes.InvalidArgument for missing IDs, or codes.Internal when the underlying connect operation fails.
func (s *Service) ConnectNetwork(
	ctx context.Context,
	req *protos.ConnectNetworkRequest,
) (*protos.ConnectNetworkResponse, error) {
	if req.GetNetworkId() == "" {
		return nil, status.Error(codes.InvalidArgument, "network ID is required")
	}

	if req.GetContainerId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	err := s.layer.ConnectNetwork(ctx, req.GetNetworkId(), req.GetContainerId(), mapEndpointSettings(req.GetSettings()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot connect network: %v", err)
	}

	return &protos.ConnectNetworkResponse{}, nil
}
