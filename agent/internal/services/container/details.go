// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package container

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetContainerDetails(
	ctx context.Context,
	req *protos.GetContainerDetailsRequest,
) (*protos.GetContainerDetailsResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	details, err := s.layer.GetContainerDetails(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get container details: %v", err)
	}

	return &protos.GetContainerDetailsResponse{
		Id:       details.ID,
		Image:    details.Image,
		Name:     details.Name,
		Status:   extractStatus(details.State),
		Created:  details.Created,
		Mounts:   mapMounts(details.Mounts),
		Networks: mapNetworks(details.NetworkSettings.Networks),
	}, nil
}
