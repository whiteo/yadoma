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

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
)

// GetVolumes lists Docker volumes using the provided volume.ListOptions.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration and returns a volume.ListResponse on success.
// On failure, it returns an error wrapped with additional context information.
func (l *Layer) GetVolumes(ctx context.Context, opts volume.ListOptions) (volume.ListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	volumes, err := l.client.VolumeList(ctx, opts)
	if err != nil {
		return volume.ListResponse{}, fmt.Errorf("cannot list volume: %w", err)
	}
	return volumes, nil
}

// GetVolumeDetails inspects a Docker volume by its name \(identifier\).
// It derives a context with a predefined timeout \(ctxTimeout\) from the incoming context
// to bound the operation duration and returns the full volume metadata on success.
// On failure, it returns an error wrapped with additional context information, including the volume identifier.
func (l *Layer) GetVolumeDetails(ctx context.Context, id string) (volume.Volume, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	vol, err := l.client.VolumeInspect(ctx, id)
	if err != nil {
		return volume.Volume{}, fmt.Errorf("cannot inspect volume %s: %w", id, err)
	}
	return vol, nil
}

// CreateVolume creates a Docker volume using the provided volume.CreateOptions.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration and returns the created volume metadata on success.
// On failure, it returns an error wrapped with additional context information.
func (l *Layer) CreateVolume(ctx context.Context, opts volume.CreateOptions) (volume.Volume, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	vol, err := l.client.VolumeCreate(ctx, opts)
	if err != nil {
		return volume.Volume{}, fmt.Errorf("cannot create volume: %w", err)
	}
	return vol, nil
}

// RemoveVolume removes a Docker volume by its name (identifier).
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration. If force is true, the volume is removed even if it is in use.
// On failure, it returns an error wrapped with additional context information, including the volume identifier.
func (l *Layer) RemoveVolume(ctx context.Context, id string, force bool) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := l.client.VolumeRemove(ctx, id, force); err != nil {
		return fmt.Errorf("cannot remove volume %s: %w", id, err)
	}
	return nil
}

// PruneVolumes removes unused Docker volumes according to the provided filter arguments.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration and returns a volume.PruneReport on success.
// On failure, it returns an error wrapped with additional context information.
func (l *Layer) PruneVolumes(ctx context.Context, args filters.Args) (volume.PruneReport, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	vol, err := l.client.VolumesPrune(ctx, args)
	if err != nil {
		return volume.PruneReport{}, fmt.Errorf("cannot prune volumes: %w", err)
	}
	return vol, nil
}
