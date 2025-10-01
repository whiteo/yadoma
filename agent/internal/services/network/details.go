// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package network

import (
	"context"
	"time"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/network"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
