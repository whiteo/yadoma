// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package image

import (
	"context"
	"io"

	docker "github.com/whiteo/yadoma/internal/dockers"
	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"

	"google.golang.org/grpc"
)

type layerAPI interface {
	GetImages(ctx context.Context, opts image.ListOptions) ([]image.Summary, error)
	GetImageDetails(ctx context.Context, id string) (image.InspectResponse, error)
	PullImage(ctx context.Context, imageName string, opts image.PullOptions) (io.ReadCloser, error)
	RemoveImage(ctx context.Context, imageID string, opts image.RemoveOptions) ([]image.DeleteResponse, error)
	PruneImage(ctx context.Context, args filters.Args) (image.PruneReport, error)
	BuildImage(
		ctx context.Context,
		buildContext io.Reader,
		opts build.ImageBuildOptions,
	) (build.ImageBuildResponse, error)
}

type Service struct {
	protos.UnimplementedImageServiceServer
	layer layerAPI
}

func NewImageService(layer *docker.Layer) *Service {
	return &Service{layer: layer}
}

func (s *Service) Register(rpc *grpc.Server) {
	protos.RegisterImageServiceServer(rpc, s)
}
