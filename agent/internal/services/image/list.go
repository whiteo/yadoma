// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package image

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/image"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
