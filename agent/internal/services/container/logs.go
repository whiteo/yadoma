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
