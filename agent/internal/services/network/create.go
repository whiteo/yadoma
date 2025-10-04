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

// CreateNetwork creates a Docker network using the provided name and options mapped from the request.
// It validates that the network name is non-empty and uses the incoming context for cancellation and deadlines.
// On success, it returns the created network ID.
// On failure, it returns a gRPC error (codes.InvalidArgument for bad input, codes.Internal for Docker layer errors).
func (s *Service) CreateNetwork(
	ctx context.Context,
	req *protos.CreateNetworkRequest,
) (*protos.CreateNetworkResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "network name is required")
	}

	r, err := s.layer.CreateNetwork(ctx, req.Name, mapCreateOptions(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create network: %v", err)
	}

	return &protos.CreateNetworkResponse{Id: r.ID}, nil
}
