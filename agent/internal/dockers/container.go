// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func (l *Layer) GetContainers(ctx context.Context, opts container.ListOptions) ([]container.Summary, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	containers, err := l.client.ContainerList(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot list containers: %w", err)
	}
	return containers, nil
}

func (l *Layer) GetContainerDetails(ctx context.Context, id string) (container.InspectResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	c, err := l.client.ContainerInspect(ctx, id)
	if err != nil {
		return container.InspectResponse{}, fmt.Errorf("cannot inspect container %s: %w", id, err)
	}
	return c, nil
}

func (l *Layer) GetContainerLogs(ctx context.Context, id string, opts container.LogsOptions) (io.ReadCloser, error) {
	// Не создаем timeout контекст для стриминговых операций логов,
	// так как они могут выполняться долго при follow=true
	// Используем оригинальный контекст от клиента
	logs, err := l.client.ContainerLogs(ctx, id, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot get logs for container %s: %w", id, err)
	}
	return logs, nil
}

func (l *Layer) GetContainerStats(ctx context.Context,
	id string,
	stream bool,
) (container.StatsResponseReader, error) {
	// Не создаем timeout контекст для стриминговых операций,
	// так как они могут выполняться долго
	// Используем оригинальный контекст от клиента
	stats, err := l.client.ContainerStats(ctx, id, stream)
	if err != nil {
		return container.StatsResponseReader{}, fmt.Errorf("cannot get stats for container %s: %w", id, err)
	}
	return stats, nil
}

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

func (l *Layer) RemoveContainer(ctx context.Context, id string, opts container.RemoveOptions) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := l.client.ContainerRemove(ctx, id, opts); err != nil {
		return fmt.Errorf("cannot remove container %s: %w", id, err)
	}
	return nil
}

func (l *Layer) StartContainer(ctx context.Context, id string, opts container.StartOptions) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := l.client.ContainerStart(ctx, id, opts); err != nil {
		return fmt.Errorf("cannot start container %s: %w", id, err)
	}
	return nil
}

func (l *Layer) StopContainer(ctx context.Context, id string, opts container.StopOptions) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := l.client.ContainerStop(ctx, id, opts); err != nil {
		return fmt.Errorf("cannot stop container %s: %w", id, err)
	}
	return nil
}

func (l *Layer) RestartContainer(ctx context.Context, id string, opts container.StopOptions) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := l.client.ContainerRestart(ctx, id, opts); err != nil {
		return fmt.Errorf("cannot restart container %s: %w", id, err)
	}
	return nil
}

func (l *Layer) PauseContainer(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := l.client.ContainerPause(ctx, id); err != nil {
		return fmt.Errorf("cannot pause container %s: %w", id, err)
	}
	return nil
}

func (l *Layer) UnpauseContainer(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := l.client.ContainerUnpause(ctx, id); err != nil {
		return fmt.Errorf("cannot unpause container %s: %w", id, err)
	}
	return nil
}

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

func (l *Layer) RenameContainer(ctx context.Context, id, name string) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := l.client.ContainerRename(ctx, id, name); err != nil {
		return fmt.Errorf("cannot rename container %s to %s: %w", id, name, err)
	}
	return nil
}
