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
	"github.com/whiteo/yadoma/internal/protos"
	service "github.com/whiteo/yadoma/internal/services"

	"github.com/docker/docker/api/types/container"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetContainerLogs streams logs of a Docker container to the gRPC client.
// It validates the request, constructs Docker `container.LogsOptions` (currently supports `Follow`),
// and acquires a log reader from the Docker layer using the incoming context for cancellation.
// Log data is forwarded to the client as chunked `GetContainerLogsResponse` messages until EOF or context cancellation.
// Returns `InvalidArgument` for an empty container ID and `Internal` for failures obtaining or streaming logs.
func (s *Service) GetContainerLogs(
	req *protos.GetContainerLogsRequest,
	stream protos.ContainerService_GetContainerLogsServer,
) error {
	if req.GetId() == "" {
		return status.Error(codes.InvalidArgument, "container ID is required")
	}

	opts := container.LogsOptions{
		Follow: req.GetFollow(),
	}

	logsReader, err := s.layer.GetContainerLogs(stream.Context(), req.GetId(), opts)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot get container logs: %v", err)
	}
	defer func() {
		if cErr := logsReader.Close(); cErr != nil {
			log.Error().Err(cErr).Msg("error closing logs reader")
		}
	}()

	return service.StreamReader(logsReader, func(chunk []byte) error {
		return stream.Send(&protos.GetContainerLogsResponse{Chunk: chunk})
	})
}
