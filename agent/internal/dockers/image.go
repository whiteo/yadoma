// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
)

func (l *Layer) GetImages(ctx context.Context, opts image.ListOptions) ([]image.Summary, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	images, err := l.client.ImageList(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot list image: %w", err)
	}
	return images, nil
}

func (l *Layer) GetImageDetails(ctx context.Context, id string) (image.InspectResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	img, err := l.client.ImageInspect(ctx, id)
	if err != nil {
		return image.InspectResponse{}, fmt.Errorf("cannot inspect image %s: %w", id, err)
	}
	return img, nil
}

func (l *Layer) RemoveImage(ctx context.Context,
	id string,
	opts image.RemoveOptions,
) ([]image.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	img, err := l.client.ImageRemove(ctx, id, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot remove image %s: %w", id, err)
	}
	return img, nil
}

func (l *Layer) PullImage(ctx context.Context,
	link string,
	opts image.PullOptions,
) (io.ReadCloser, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	pull, err := l.client.ImagePull(ctx, link, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot pull image %s: %w", link, err)
	}
	return pull, nil
}

func (l *Layer) BuildImage(ctx context.Context,
	buildCtx io.Reader,
	opts build.ImageBuildOptions,
) (build.ImageBuildResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	resp, err := l.client.ImageBuild(ctx, buildCtx, opts)
	if err != nil {
		return build.ImageBuildResponse{}, fmt.Errorf("cannot build image: %w", err)
	}
	return resp, nil
}

func (l *Layer) PruneImage(ctx context.Context, args filters.Args) (image.PruneReport, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	report, err := l.client.ImagesPrune(ctx, args)
	if err != nil {
		return image.PruneReport{}, fmt.Errorf("cannot prune image: %w", err)
	}
	return report, nil
}
