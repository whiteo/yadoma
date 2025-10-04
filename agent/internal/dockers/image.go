// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package docker provides a thin, internal wrapper over the Docker Engine API client.
// It centralizes container, image, network, volume, and system operations while keeping
// calls close to the upstream API. Most requests are bounded with a package-level timeout
// (ctxTimeout) derived from the caller's context to prevent indefinite waits.
//
// Streaming endpoints (for example, logs and stats) use the caller's context as-is.
// Callers must read from and close returned streams. The package does not spawn
// goroutines on behalf of the caller and relies on context cancellation for shutdown.
//
// Errors are returned with additional context to aid diagnostics. Configuration,
// retries, and higher-level policies are left to callers. The package is intended
// for internal use by services that compose these primitives.
package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
)

// GetImages lists Docker images using the provided image.ListOptions.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration and returns image summaries on success.
// On failure, it returns an error wrapped with additional context information.
func (l *Layer) GetImages(ctx context.Context, opts image.ListOptions) ([]image.Summary, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	images, err := l.client.ImageList(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot list image: %w", err)
	}
	return images, nil
}

// GetImageDetails retrieves detailed information about a Docker image by its ID.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration and returns an image.InspectResponse on success.
// On failure, it returns an error wrapped with additional context information, including the image ID.
func (l *Layer) GetImageDetails(ctx context.Context, id string) (image.InspectResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	img, err := l.client.ImageInspect(ctx, id)
	if err != nil {
		return image.InspectResponse{}, fmt.Errorf("cannot inspect image %s: %w", id, err)
	}
	return img, nil
}

// RemoveImage deletes a Docker image by its ID using the provided image.RemoveOptions.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation and returns a slice of image.DeleteResponse entries on success.
// On failure, it returns an error wrapped with additional context information, including the image ID.
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

// PullImage downloads a Docker image by reference using the provided image.PullOptions.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration and returns a stream (io.ReadCloser) on success.
// The caller must read from and close the returned stream.
// On failure, it returns an error wrapped with additional context information, including the image reference.
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

// BuildImage builds a Docker image from the provided build context using the given build.ImageBuildOptions.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration and returns a build.ImageBuildResponse on success.
// The response contains a streaming body; callers must read from and close resp.Body.
// On failure, it returns an error wrapped with additional context information.
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

// PruneImage removes unused Docker images matching the provided filters.Args.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration and returns an image.PruneReport on success.
// On failure, it returns an error wrapped with additional context information.
func (l *Layer) PruneImage(ctx context.Context, args filters.Args) (image.PruneReport, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	report, err := l.client.ImagesPrune(ctx, args)
	if err != nil {
		return image.PruneReport{}, fmt.Errorf("cannot prune image: %w", err)
	}
	return report, nil
}
