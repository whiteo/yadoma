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

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// GetContainers lists Docker containers using the provided container.ListOptions.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration and returns container summaries on success.
// On failure, it returns an error wrapped with additional context information.
func (l *Layer) GetContainers(ctx context.Context, opts container.ListOptions) ([]container.Summary, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	containers, err := l.client.ContainerList(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot list containers: %w", err)
	}
	return containers, nil
}

// GetContainerDetails inspects a Docker container by its ID.
// It uses a context derived with the predefined ctxTimeout and returns
// the full container.InspectResponse on success, or a wrapped error on failure.
func (l *Layer) GetContainerDetails(ctx context.Context, id string) (container.InspectResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	c, err := l.client.ContainerInspect(ctx, id)
	if err != nil {
		return container.InspectResponse{}, fmt.Errorf("cannot inspect container %s: %w", id, err)
	}
	return c, nil
}

// GetContainerLogs returns a streaming reader for the logs of the container identified by id.
// It invokes the Docker API with the provided container.LogsOptions and returns an io.ReadCloser.
// The caller must read from and close the returned stream to avoid leaks.
// If opts.Follow is true, the stream stays open until ctx is canceled or the log stream ends.
// For TTY-enabled containers, logs are not multiplexed; otherwise stdout/stderr are multiplexed per the Docker API.
// The provided ctx is used as-is (no internal timeout).
// On failure, it returns a wrapped error with context.
func (l *Layer) GetContainerLogs(ctx context.Context, id string, opts container.LogsOptions) (io.ReadCloser, error) {
	logs, err := l.client.ContainerLogs(ctx, id, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot get logs for container %s: %w", id, err)
	}
	return logs, nil
}

// GetContainerStats retrieves live resource-usage statistics for the container identified by id.
// If stream is true, the Docker daemon keeps the connection open and continuously sends stats
// until ctx is canceled; otherwise, a single snapshot is returned.
// The provided ctx is used as-is. Cancel it to stop a streaming request.
// On success, it returns a StatsResponseReader that the caller must read from and close.
// On failure, it returns an error propagated from the Docker client.
func (l *Layer) GetContainerStats(ctx context.Context,
	id string,
	stream bool,
) (container.StatsResponseReader, error) {
	stats, err := l.client.ContainerStats(ctx, id, stream)
	if err != nil {
		return container.StatsResponseReader{}, fmt.Errorf("cannot get stats for container %s: %w", id, err)
	}
	return stats, nil
}

// CreateContainer creates a Docker container using the provided specifications.
// A child context with a predefined timeout (ctxTimeout) is derived from ctx
// to bound the request duration.
// Note: This method only creates the container; it does not start it.
// On success, it returns container.CreateResponse. On failure, it returns an
// error wrapped with context about the container name.
func (l *Layer) CreateContainer(
	ctx context.Context,
	config *container.Config,
	hostConfig *container.HostConfig,
	networkingConfig *network.NetworkingConfig,
	platform *ocispec.Platform,
	containerName string,
) (container.CreateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	resp, err := l.client.ContainerCreate(
		ctx,
		config,
		hostConfig,
		networkingConfig,
		platform,
		containerName,
	)
	if err != nil {
		return container.CreateResponse{}, fmt.Errorf("cannot create container %s: %w", containerName, err)
	}

	return resp, nil
}

// RemoveContainer removes a Docker container identified by id using the provided container.RemoveOptions.
// A child context with the predefined timeout (ctxTimeout) is derived from ctx to bound the request duration.
// Depending on opts, this may force-remove a running container and/or delete its associated volumes.
// Returns nil on success; on failure, returns an error wrapped with the container id and underlying cause.
func (l *Layer) RemoveContainer(ctx context.Context, id string, opts container.RemoveOptions) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := l.client.ContainerRemove(ctx, id, opts); err != nil {
		return fmt.Errorf("cannot remove container %s: %w", id, err)
	}
	return nil
}

// StartContainer starts an existing Docker container identified by id.
// It derives a child context with the predefined timeout (ctxTimeout) from ctx
// to bound the operation duration.
// The call delegates to Docker's ContainerStart API and applies the provided
// container.StartOptions.
// Returns nil on success; on failure, returns an error wrapped with the
// container id and underlying cause.
func (l *Layer) StartContainer(ctx context.Context, id string, opts container.StartOptions) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := l.client.ContainerStart(ctx, id, opts); err != nil {
		return fmt.Errorf("cannot start container %s: %w", id, err)
	}
	return nil
}

// StopContainer gracefully stops a Docker container identified by id.
// A child context with the predefined ctxTimeout is derived from ctx to bound the operation.
// The call delegates to Docker's ContainerStop using the provided container.StopOptions.
// Returns nil on success; on failure, returns an error wrapped with the container id.
// Cancel ctx to abort the request early.
func (l *Layer) StopContainer(ctx context.Context, id string, opts container.StopOptions) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := l.client.ContainerStop(ctx, id, opts); err != nil {
		return fmt.Errorf("cannot stop container %s: %w", id, err)
	}
	return nil
}

// RestartContainer restarts a Docker container identified by id.
// A child context with the predefined timeout (ctxTimeout) is derived from ctx
// to bound the operation duration.
// It delegates to Docker's ContainerRestart and accepts container.StopOptions
// to control the grace period before the container is forcibly terminated.
// Returns nil on success; on failure, returns an error wrapped with the
// container id and underlying cause.
func (l *Layer) RestartContainer(ctx context.Context, id string, opts container.StopOptions) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := l.client.ContainerRestart(ctx, id, opts); err != nil {
		return fmt.Errorf("cannot restart container %s: %w", id, err)
	}
	return nil
}

// PauseContainer pauses a Docker container identified by id.
// It derives a child context with the predefined timeout (ctxTimeout) to bound the operation.
// The call delegates to Docker\'s ContainerPause API.
// Returns nil on success; on failure, returns an error wrapped with the container id.
func (l *Layer) PauseContainer(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := l.client.ContainerPause(ctx, id); err != nil {
		return fmt.Errorf("cannot pause container %s: %w", id, err)
	}
	return nil
}

// UnpauseContainer resumes a previously paused Docker container identified by id.
// It derives a child context with the predefined timeout (ctxTimeout) to bound the operation.
// The call delegates to Docker's ContainerUnpause API.
// Returns nil on success; on failure, returns an error wrapped with the container id.
func (l *Layer) UnpauseContainer(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := l.client.ContainerUnpause(ctx, id); err != nil {
		return fmt.Errorf("cannot unpause container %s: %w", id, err)
	}
	return nil
}

// KillContainer sends a kill signal to the Docker container identified by id.
// It derives a child context with the predefined timeout (ctxTimeout) from ctx
// to bound the operation duration.
// If signal is empty, "SIGKILL" is used by default.
// The signal can be specified as a POSIX name (e.g., "SIGTERM") or a number (e.g., "9"),
// as supported by the Docker daemon.
// Returns nil on success; on failure, returns an error wrapped with the container id and signal.
func (l *Layer) KillContainer(ctx context.Context, id, signal string) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if signal == "" {
		signal = "SIGKILL"
	}

	if err := l.client.ContainerKill(ctx, id, signal); err != nil {
		return fmt.Errorf("cannot kill container %s with signal %s: %w", id, signal, err)
	}
	return nil
}

// RenameContainer renames an existing Docker container identified by id.
// A child context with the predefined timeout (ctxTimeout) is derived from ctx
// to bound the operation duration.
// The call delegates to Docker's ContainerRename API.
// The new name must be unique and conform to Docker's container-naming rules.
// Returns nil on success; on failure, returns an error wrapped with the container id,
// target name, and the underlying cause.
func (l *Layer) RenameContainer(ctx context.Context, id, name string) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := l.client.ContainerRename(ctx, id, name); err != nil {
		return fmt.Errorf("cannot rename container %s to %s: %w", id, name, err)
	}
	return nil
}
