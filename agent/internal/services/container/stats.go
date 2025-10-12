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
	containerID := req.GetId()
	if containerID == "" {
		return status.Error(codes.InvalidArgument, "container ID is required")
	}

	log.Debug().
		Str("container", containerID).
		Bool("stream", req.GetStream()).
		Msg("Starting stats stream")

	statsReader, err := s.layer.GetContainerStats(stream.Context(), containerID, req.GetStream())
	if err != nil {
		log.Error().
			Err(err).
			Str("container", containerID).
			Msg("Failed to get stats reader from Docker")
		return status.Errorf(codes.Internal, "cannot get container stats: %v", err)
	}
	defer func() {
		if cErr := statsReader.Body.Close(); cErr != nil {
			log.Error().Err(cErr).Str("container", containerID).Msg("error closing stats reader")
		} else {
			log.Debug().Str("container", containerID).Msg("Stats reader closed")
		}
	}()

	log.Debug().Str("container", containerID).Msg("Starting to decode stats stream")

	err = service.StreamDecoder[container.StatsResponse](statsReader.Body, func(stats container.StatsResponse) error {
		log.Trace().
			Str("container", containerID).
			Uint64("cpu", stats.CPUStats.CPUUsage.TotalUsage).
			Uint64("mem", stats.MemoryStats.Usage).
			Msg("Decoded stats, sending to client")
		return stream.Send(mapStats(stats))
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("container", containerID).
			Msg("Error in stats stream decoder")
	} else {
		log.Debug().Str("container", containerID).Msg("Stats stream completed successfully")
	}

	return err
}
