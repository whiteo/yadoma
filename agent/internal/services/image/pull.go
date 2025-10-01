// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package image

import (
	"github.com/whiteo/yadoma/internal/protos"
	service "github.com/whiteo/yadoma/internal/services"

	"github.com/docker/docker/api/types/image"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) PullImage(req *protos.PullImageRequest,
	stream protos.ImageService_PullImageServer,
) error {
	if req.GetLink() == "" {
		return status.Error(codes.InvalidArgument, "image ID is required")
	}

	pullReader, err := s.layer.PullImage(stream.Context(), req.Link, image.PullOptions{RegistryAuth: req.RegistryAuth})
	if err != nil {
		return status.Errorf(codes.Internal, "cannot pull image: %v", err)
	}
	defer func() {
		if cErr := pullReader.Close(); cErr != nil {
			log.Error().Err(cErr).Msg("error closing pull reader")
		}
	}()

	return service.StreamReader(pullReader, func(chunk []byte) error {
		return stream.Send(&protos.PullImageResponse{Chunk: chunk})
	})
}
