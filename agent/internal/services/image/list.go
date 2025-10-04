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

	"github.com/docker/docker/api/types/image"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetImages lists Docker images based on the request parameters.
// It delegates to the Docker layer with image.ListOptions populated from the request,
// maps the results into protobuf messages, and returns a protos.GetImagesResponse.
// Errors from the Docker layer are translated into a gRPC codes.Internal status.
// The call respects the incoming context for cancellation and deadlines.
func (s *Service) GetImages(ctx context.Context, req *protos.GetImagesRequest) (*protos.GetImagesResponse, error) {
	list, err := s.layer.GetImages(ctx, image.ListOptions{All: req.GetAll()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot list images: %v", err)
	}

	resp := &protos.GetImagesResponse{
		Images: make([]*protos.GetImageResponse, 0, len(list)),
	}

	for _, i := range list {
		resp.Images = append(resp.Images, &protos.GetImageResponse{
			Id:         i.ID,
			RepoTags:   i.RepoTags,
			Created:    i.Created,
			Size:       i.Size,
			Containers: i.Containers,
		})
	}

	return resp, nil
}
