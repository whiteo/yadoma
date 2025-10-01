// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/system"
)

func (l *Layer) GetSystemInfo(ctx context.Context) (system.Info, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	info, err := l.client.Info(ctx)
	if err != nil {
		return system.Info{}, fmt.Errorf("cannot get system info: %w", err)
	}
	return info, nil
}

func (l *Layer) GetDiskUsage(ctx context.Context, opts types.DiskUsageOptions) (types.DiskUsage, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	du, err := l.client.DiskUsage(ctx, opts)
	if err != nil {
		return types.DiskUsage{}, fmt.Errorf("cannot get disk usage: %w", err)
	}
	return du, nil
}
