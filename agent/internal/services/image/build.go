// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package image

import (
	"context"
	"time"

	"github.com/whiteo/yadoma/internal/protos"
	service "github.com/whiteo/yadoma/internal/services"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) BuildImage(
	req *protos.BuildImageRequest,
	stream protos.ImageService_BuildImageServer,
) error {
	if req.Dockerfile == "" {
		return status.Error(codes.InvalidArgument, "dockerfile is required")
	}

	ctx, cancel := context.WithTimeout(stream.Context(), 30*time.Second)
	defer cancel()

	opts, buildCtx := mapBuildOptions(req)

	resp, err := s.layer.BuildImage(ctx, buildCtx, opts)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			log.Error().Err(cErr).Msg("error closing build reader")
		}
	}()

	return service.StreamReader(resp.Body, func(chunk []byte) error {
		return stream.Send(&protos.BuildImageResponse{Chunk: chunk})
	})
}
