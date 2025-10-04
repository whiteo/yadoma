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
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/system"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

const ctxTimeout = 30 * time.Second

type ClientInterface interface {
	// Container methods
	ContainerList(ctx context.Context, options container.ListOptions) ([]container.Summary, error)
	ContainerInspect(ctx context.Context, containerID string) (container.InspectResponse, error)
	ContainerLogs(ctx context.Context, containerID string, options container.LogsOptions) (io.ReadCloser, error)
	ContainerStats(ctx context.Context, containerID string, stream bool) (container.StatsResponseReader, error)
	ContainerCreate(ctx context.Context,
		config *container.Config,
		hostConfig *container.HostConfig,
		networkingConfig *network.NetworkingConfig,
		platform *ocispec.Platform,
		containerName string,
	) (container.CreateResponse, error)
	ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error
	ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error
	ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerPause(ctx context.Context, containerID string) error
	ContainerUnpause(ctx context.Context, containerID string) error
	ContainerKill(ctx context.Context, containerID string, signal string) error
	ContainerRename(ctx context.Context, containerID string, newName string) error

	// Image methods
	ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
	ImageInspect(ctx context.Context,
		imageID string,
		inspectOpts ...client.ImageInspectOption,
	) (image.InspectResponse, error)
	ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error)
	ImagePull(ctx context.Context, refStr string, options image.PullOptions) (io.ReadCloser, error)
	ImageBuild(ctx context.Context,
		buildContext io.Reader,
		options build.ImageBuildOptions,
	) (build.ImageBuildResponse, error)
	ImagesPrune(ctx context.Context, pruneFilters filters.Args) (image.PruneReport, error)

	// Network methods
	NetworkList(ctx context.Context, options network.ListOptions) ([]network.Summary, error)
	NetworkInspect(ctx context.Context, networkID string, options network.InspectOptions) (network.Inspect, error)
	NetworkCreate(ctx context.Context, name string, options network.CreateOptions) (network.CreateResponse, error)
	NetworkConnect(ctx context.Context, networkID, containerID string, config *network.EndpointSettings) error
	NetworkDisconnect(ctx context.Context, networkID, containerID string, force bool) error
	NetworkRemove(ctx context.Context, networkID string) error
	NetworksPrune(ctx context.Context, pruneFilters filters.Args) (network.PruneReport, error)

	// Volume methods
	VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error)
	VolumeInspect(ctx context.Context, volumeID string) (volume.Volume, error)
	VolumeCreate(ctx context.Context, options volume.CreateOptions) (volume.Volume, error)
	VolumeRemove(ctx context.Context, volumeID string, force bool) error
	VolumesPrune(ctx context.Context, pruneFilters filters.Args) (volume.PruneReport, error)

	// System methods
	Info(ctx context.Context) (system.Info, error)
	DiskUsage(ctx context.Context, options types.DiskUsageOptions) (types.DiskUsage, error)
}

type Layer struct {
	client ClientInterface
}

// NewLayer constructs a Layer that wraps the provided Docker Engine API client.
// It binds the given client to enable container, image, network, volume, and
// system operations through this package.
// Non-streaming requests are bounded by the package-level timeout (ctxTimeout);
// streaming endpoints use the caller's context.
// The caller retains ownership of the client and should close it when finished.
func NewLayer(c *client.Client) *Layer {
	return &Layer{client: c}
}
