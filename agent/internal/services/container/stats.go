// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package container

import (
	"github.com/whiteo/yadoma/internal/protos"
	service "github.com/whiteo/yadoma/internal/services"

	"github.com/docker/docker/api/types/container"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
