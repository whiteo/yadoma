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

	"github.com/whiteo/yadoma/internal/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetImageDetails retrieves detailed metadata for a Docker image by its ID.
// It validates that the request contains a non-empty image ID, delegates the lookup
// to the Docker layer, and maps the result into the protobuf response type.
// On error, it returns gRPC status errors: InvalidArgument for a missing ID and
// Internal for failures returned by the Docker layer.
func (s *Service) GetImageDetails(
	ctx context.Context,
	req *protos.GetImageDetailsRequest,
) (*protos.GetImageDetailsResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "image ID is required")
	}

	details, err := s.layer.GetImageDetails(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get image details: %v", err)
	}

	return &protos.GetImageDetailsResponse{
		Id:           details.ID,
		RepoTags:     details.RepoTags,
		Created:      details.Created,
		Size:         details.Size,
		Author:       details.Author,
		Architecture: details.Architecture,
		Os:           details.Os,
	}, nil
}
