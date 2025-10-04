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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/system"
)

// GetSystemInfo retrieves Docker Engine system information.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration and returns the engine's system info on success.
// On failure, it returns an error wrapped with additional context information.
func (l *Layer) GetSystemInfo(ctx context.Context) (system.Info, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	info, err := l.client.Info(ctx)
	if err != nil {
		return system.Info{}, fmt.Errorf("cannot get system info: %w", err)
	}
	return info, nil
}

// GetDiskUsage queries Docker Engine disk usage using the provided types.DiskUsageOptions.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration and returns aggregated disk usage details on success.
// On failure, it returns an error wrapped with additional context information.
func (l *Layer) GetDiskUsage(ctx context.Context, opts types.DiskUsageOptions) (types.DiskUsage, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	du, err := l.client.DiskUsage(ctx, opts)
	if err != nil {
		return types.DiskUsage{}, fmt.Errorf("cannot get disk usage: %w", err)
	}
	return du, nil
}
