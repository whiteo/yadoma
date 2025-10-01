// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package image

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
