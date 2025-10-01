// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package image

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/filters"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
