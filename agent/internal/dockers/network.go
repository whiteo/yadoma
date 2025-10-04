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
	"github.com/docker/docker/api/types/network"
)

// GetNetworks lists Docker networks using the provided network.ListOptions.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration and returns network summaries on success.
// On failure, it returns an error wrapped with additional context information.
func (l *Layer) GetNetworks(ctx context.Context, opts network.ListOptions) ([]network.Summary, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	networks, err := l.client.NetworkList(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot list network: %w", err)
	}
	return networks, nil
}

// GetNetworkDetails inspects a Docker network by ID using the provided network.InspectOptions.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration and returns a detailed network inspection on success.
// On failure, it returns a zero-value network.Inspect and an error wrapped with additional context.
func (l *Layer) GetNetworkDetails(ctx context.Context,
	id string,
	opts network.InspectOptions,
) (network.Inspect, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	n, err := l.client.NetworkInspect(ctx, id, opts)
	if err != nil {
		return network.Inspect{}, fmt.Errorf("cannot inspect network %s: %w", id, err)
	}
	return n, nil
}

// CreateNetwork creates a Docker network using the provided network.CreateOptions.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration and returns a network.CreateResponse on success.
// On failure, it returns a zero-value network.CreateResponse and an error wrapped with additional context.
func (l *Layer) CreateNetwork(ctx context.Context,
	name string,
	opts network.CreateOptions,
) (network.CreateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	n, err := l.client.NetworkCreate(ctx, name, opts)
	if err != nil {
		return network.CreateResponse{}, fmt.Errorf("cannot create network %s: %w", name, err)
	}
	return n, nil
}

// ConnectNetwork connects a Docker container to a network using the provided network.EndpointSettings.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration. On success, it returns nil.
// On failure, it returns an error wrapped with additional context that includes the container and network identifiers.
func (l *Layer) ConnectNetwork(ctx context.Context,
	networkID,
	containerID string,
	config *network.EndpointSettings,
) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	err := l.client.NetworkConnect(ctx, networkID, containerID, config)
	if err != nil {
		return fmt.Errorf("cannot connect container %s to network %s: %w", containerID, networkID, err)
	}
	return nil
}

// DisconnectNetwork disconnects a Docker container from a network.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration.
// If force is true, the container is forcibly disconnected.
// On success, it returns nil.
// On failure, it returns an error wrapped with additional context that includes the container and network identifiers.
func (l *Layer) DisconnectNetwork(ctx context.Context, networkID, containerID string, force bool) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	err := l.client.NetworkDisconnect(ctx, networkID, containerID, force)
	if err != nil {
		return fmt.Errorf("cannot disconnect container %s from network %s: %w", containerID, networkID, err)
	}
	return nil
}

// PruneNetworks removes unused Docker networks matching the provided filters.Args.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration and returns a network.PruneReport on success.
// On failure, it returns a zero-value network.PruneReport and an error wrapped with additional context.
func (l *Layer) PruneNetworks(ctx context.Context, args filters.Args) (network.PruneReport, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	report, err := l.client.NetworksPrune(ctx, args)
	if err != nil {
		return network.PruneReport{}, fmt.Errorf("cannot prune networks: %w", err)
	}
	return report, nil
}

// RemoveNetwork removes a Docker network by ID.
// It derives a context with a predefined timeout (ctxTimeout) from the incoming context
// to bound the operation duration.
// On success, it returns nil.
// On failure, it returns an error wrapped with additional context including the network ID.
func (l *Layer) RemoveNetwork(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	err := l.client.NetworkRemove(ctx, id)
	if err != nil {
		return fmt.Errorf("cannot remove network %s: %w", id, err)
	}
	return nil
}
