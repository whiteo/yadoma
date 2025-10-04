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

	"github.com/docker/docker/api/types/filters"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PruneImages prunes unused Docker images based on the request options.
// It forwards the incoming context and applies a Docker filter with the 'All' flag
// to determine the pruning scope. On success, it returns a response containing
// the list of deleted/untagged images and the total reclaimed space. On failure,
// it returns a gRPC error with code Internal.
func (s *Service) PruneImages(
	ctx context.Context,
	req *protos.PruneImagesRequest,
) (*protos.PruneImagesResponse, error) {
	count, err := s.layer.PruneImage(ctx, filters.NewArgs(filters.Arg("All", req.All)))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot prune image: %v", err)
	}

	resp := &protos.PruneImagesResponse{
		ImagesDeleted:  make([]*protos.PruneImageResult, 0, len(count.ImagesDeleted)),
		SpaceReclaimed: count.SpaceReclaimed,
	}

	for _, id := range count.ImagesDeleted {
		resp.ImagesDeleted = append(resp.ImagesDeleted, &protos.PruneImageResult{
			Deleted:  id.Deleted,
			Untagged: id.Untagged,
		})
	}

	return resp, nil
}
