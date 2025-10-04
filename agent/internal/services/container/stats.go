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

// GetContainerStats streams Docker container statistics via gRPC.
// Validates that a container ID is provided, then requests stats from the Docker layer,
// honoring the incoming context for cancellation and closing the response body on exit.
// Decodes each Docker stats payload and sends it to the client after mapping to protobuf.
// Returns gRPC errors: InvalidArgument for a missing ID, Internal on Docker access failures;
// supports single-shot or continuous updates based on the request's stream flag.
func (s *Service) GetContainerStats(
	req *protos.GetContainerStatsRequest,
	stream protos.ContainerService_GetContainerStatsServer,
) error {
	if req.GetId() == "" {
		return status.Error(codes.InvalidArgument, "container ID is required")
	}

	statsReader, err := s.layer.GetContainerStats(stream.Context(), req.GetId(), req.GetStream())
	if err != nil {
		return status.Errorf(codes.Internal, "cannot get container stats: %v", err)
	}
	defer func() {
		if cErr := statsReader.Body.Close(); cErr != nil {
			log.Error().Err(cErr).Msg("error closing stats reader")
		}
	}()

	return service.StreamDecoder[container.StatsResponse](statsReader.Body, func(stats container.StatsResponse) error {
		return stream.Send(mapStats(stats))
	})
}
