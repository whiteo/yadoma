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
	"context"
	"time"

	"github.com/whiteo/yadoma/internal/protos"
	service "github.com/whiteo/yadoma/internal/services"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BuildImage builds a Docker image and streams build output to the client.
// It validates that a Dockerfile path is provided, derives a 30s timeout from the
// incoming stream context, and delegates the build to the Docker layer.
// The build output is read incrementally and forwarded to the gRPC stream as chunks.
// The response body is closed on completion or error.
// Returns gRPC errors with appropriate codes: InvalidArgument for bad input,
// Internal for Docker or I/O failures.
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
