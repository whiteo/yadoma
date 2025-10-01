// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package container

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) CreateContainer(
	ctx context.Context,
	req *protos.CreateContainerRequest,
) (*protos.CreateContainerResponse, error) {
	if req.GetImage() == "" {
		return nil, status.Error(codes.InvalidArgument, "image is required")
	}

	config := mapConfig(req)
	hostConfig := mapHostConfig(req.GetHostConfig())
	networkingConfig := mapNetworking(req.GetNetworks())

	resp, err := s.layer.CreateContainer(ctx, config, hostConfig, networkingConfig, &ocispec.Platform{}, req.GetName())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create container: %v", err)
	}

	return &protos.CreateContainerResponse{Id: resp.ID}, nil
}
