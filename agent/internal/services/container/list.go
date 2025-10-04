// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package container provides service-layer operations for managing Docker containers.
// It implements gRPC-facing logic that validates requests, invokes the Docker layer,
// maps results to protobuf messages, and returns errors as gRPC status codes.
// Supported operations cover the container lifecycle and inspection, including create,
// list, inspect, logs and stats streaming, start/stop/restart, kill, pause/unpause,
// rename, and remove. Calls respect the caller's context; streaming endpoints propagate
// cancellation and require the caller to consume and close streams. The package is
// internal to the agent and intended to be used by higher-level gRPC servers.
package container

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/container"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetContainers lists Docker containers according to the provided request.
// It forwards the `all` and `limit` options to the Docker layer using container.ListOptions
// and maps the results into protobuf messages for the response.
// The call respects the caller's context for cancellation; on failure, it returns a gRPC
// `codes.Internal` error describing the underlying cause.
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
