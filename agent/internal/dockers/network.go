// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
)

func (l *Layer) GetNetworks(ctx context.Context, opts network.ListOptions) ([]network.Summary, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	networks, err := l.client.NetworkList(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot list network: %w", err)
	}
	return networks, nil
}

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

func (l *Layer) DisconnectNetwork(ctx context.Context, networkID, containerID string, force bool) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	err := l.client.NetworkDisconnect(ctx, networkID, containerID, force)
	if err != nil {
		return fmt.Errorf("cannot disconnect container %s from network %s: %w", containerID, networkID, err)
	}
	return nil
}

func (l *Layer) PruneNetworks(ctx context.Context, args filters.Args) (network.PruneReport, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	report, err := l.client.NetworksPrune(ctx, args)
	if err != nil {
		return network.PruneReport{}, fmt.Errorf("cannot prune networks: %w", err)
	}
	return report, nil
}

func (l *Layer) RemoveNetwork(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	err := l.client.NetworkRemove(ctx, id)
	if err != nil {
		return fmt.Errorf("cannot remove network %s: %w", id, err)
	}
	return nil
}
