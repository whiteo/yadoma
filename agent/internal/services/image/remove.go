// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package image

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/image"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
