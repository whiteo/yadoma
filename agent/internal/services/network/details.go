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
	"time"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/network"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetNetworkDetails inspects a Docker network by ID and returns its metadata.
// It validates the provided network ID, delegates to the Docker layer using the
// incoming context, and maps the result into a protobuf response (creation time
// formatted as RFC3339). On invalid input it returns codes.InvalidArgument; on
// backend failure it returns a codes.Internal gRPC status with details.
func (s *Service) GetNetworkDetails(
	ctx context.Context,
	req *protos.GetNetworkDetailsRequest,
) (*protos.GetNetworkDetailsResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "network ID is required")
	}

	details, err := s.layer.GetNetworkDetails(ctx, req.Id, network.InspectOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get network details: %v", err)
	}

	return &protos.GetNetworkDetailsResponse{
		Id:         details.ID,
		Name:       details.Name,
		Created:    details.Created.Format(time.RFC3339),
		Scope:      details.Scope,
		Driver:     details.Driver,
		Internal:   details.Internal,
		Attachable: details.Attachable,
		Ingress:    details.Ingress,
		Labels:     details.Labels,
	}, nil
}
