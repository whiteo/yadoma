// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
)

func (l *Layer) GetVolumes(ctx context.Context, opts volume.ListOptions) (volume.ListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	volumes, err := l.client.VolumeList(ctx, opts)
	if err != nil {
		return volume.ListResponse{}, fmt.Errorf("cannot list volume: %w", err)
	}
	return volumes, nil
}

func (l *Layer) GetVolumeDetails(ctx context.Context, id string) (volume.Volume, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	vol, err := l.client.VolumeInspect(ctx, id)
	if err != nil {
		return volume.Volume{}, fmt.Errorf("cannot inspect volume %s: %w", id, err)
	}
	return vol, nil
}

func (l *Layer) CreateVolume(ctx context.Context, opts volume.CreateOptions) (volume.Volume, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	vol, err := l.client.VolumeCreate(ctx, opts)
	if err != nil {
		return volume.Volume{}, fmt.Errorf("cannot create volume: %w", err)
	}
	return vol, nil
}

func (l *Layer) RemoveVolume(ctx context.Context, id string, force bool) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := l.client.VolumeRemove(ctx, id, force); err != nil {
		return fmt.Errorf("cannot remove volume %s: %w", id, err)
	}
	return nil
}

func (l *Layer) PruneVolumes(ctx context.Context, args filters.Args) (volume.PruneReport, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	vol, err := l.client.VolumesPrune(ctx, args)
	if err != nil {
		return volume.PruneReport{}, fmt.Errorf("cannot prune volumes: %w", err)
	}
	return vol, nil
}
