// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package image provides the agent's service layer for Docker image management.
// It implements gRPC-facing handlers that validate requests, delegate to the
// Docker layer, map results to protobuf messages, and translate errors into gRPC
// status codes.
//
// Supported operations include building images from a context, pulling from
// registries, listing and inspecting details, removing images, and pruning
// unused images. Streaming endpoints (for example, build and pull progress)
// propagate the caller's context; callers must consume and close returned
// streams.
//
// The package does not spawn goroutines on behalf of the caller and relies on
// context deadlines and cancellation for shutdown. It is intended for internal
// use by the agent's gRPC server layer.
package image

import (
	"github.com/whiteo/yadoma/internal/protos"
	service "github.com/whiteo/yadoma/internal/services"

	"github.com/docker/docker/api/types/image"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PullImage pulls a Docker image identified by req.Link and streams progress to the client.
// It propagates the caller's context for cancellation and deadline handling and invokes
// the Docker layer with optional registry authentication from req.RegistryAuth.
// Progress output is forwarded as raw chunks over the gRPC stream, and the underlying
// reader is always closed. On failure, returns a gRPC status error; otherwise returns nil.
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
