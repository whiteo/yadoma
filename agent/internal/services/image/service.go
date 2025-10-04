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

// NewImageService creates and returns a new Image service backed by the provided Docker layer.
// It binds the service to the given layer used by gRPC handlers.
// The service does not spawn goroutines or manage the layer's lifecycle; the caller retains ownership.
func NewImageService(layer *docker.Layer) *Service {
	return &Service{layer: layer}
}

// Register registers this Image service instance with the provided gRPC server.
// It binds the service implementation to the server so that incoming RPC calls
// to the ImageService are routed here. The function does not start the server
// or manage its lifecycle; it only performs service registration.
func (s *Service) Register(rpc *grpc.Server) {
	protos.RegisterImageServiceServer(rpc, s)
}
