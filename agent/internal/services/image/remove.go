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

// RemoveImage deletes a Docker image identified by req.Id.
// It validates input (returns InvalidArgument if the ID is empty) and delegates to the Docker layer
// using image.RemoveOptions built from req.Force and req.PruneChildren, propagating the caller's context.
// On success, it returns a RemoveImageResponse with deleted and untagged results.
// Failures are translated into gRPC Internal status errors with additional context.
func (s *Service) RemoveImage(ctx context.Context,
	req *protos.RemoveImageRequest) (*protos.RemoveImageResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "image ID is required")
	}

	removed, err := s.layer.RemoveImage(ctx, req.Id, image.RemoveOptions{
		Force:         req.Force,
		PruneChildren: req.PruneChildren,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot remove image: %v", err)
	}

	resp := &protos.RemoveImageResponse{
		Results: make([]*protos.RemoveImageResult, 0, len(removed)),
	}

	for _, r := range removed {
		resp.Results = append(resp.Results, &protos.RemoveImageResult{
			Deleted:  r.Deleted,
			Untagged: r.Untagged,
		})
	}

	return resp, nil
}
