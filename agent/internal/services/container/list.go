// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package container

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/container"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetContainers(
	ctx context.Context,
	req *protos.GetContainersRequest,
) (*protos.GetContainersResponse, error) {
	opts := container.ListOptions{
		All:   req.GetAll(),
		Limit: int(req.GetLimit()),
	}

	list, err := s.layer.GetContainers(ctx, opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot list containers: %v", err)
	}

	resp := &protos.GetContainersResponse{
		Containers: make([]*protos.GetContainerResponse, 0, len(list)),
	}

	for _, c := range list {
		resp.Containers = append(resp.Containers, &protos.GetContainerResponse{
			Id:     c.ID,
			Names:  c.Names,
			Image:  c.Image,
			State:  c.State,
			Status: c.Status,
			Ports:  mapPorts(c.Ports),
		})
	}

	return resp, nil
}
