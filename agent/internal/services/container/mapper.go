// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package container provides service-layer operations for managing Docker containers.
// It implements gRPC-facing logic that validates requests, invokes the Docker layer,
// maps results to protobuf messages, and returns errors as gRPC status codes.
// Supported operations cover the container lifecycle and inspection, including create,
// list, inspect, logs and stats streaming, start/stop/restart, kill, pause/unpause,
// rename, and remove. Calls respect the caller's context; streaming endpoints propagate
// cancellation and require the caller to consume and close streams. The package is
// internal to the agent and intended to be used by higher-level gRPC servers.
package container

import (
	"fmt"
	"strconv"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

func extractStatus(state *container.State) string {
	if state == nil {
		return ""
	}
	return state.Status
}

func mapMounts(mounts []container.MountPoint) []*protos.MountPoint {
	res := make([]*protos.MountPoint, 0, len(mounts))
	for _, m := range mounts {
		res = append(res, &protos.MountPoint{
			Source:      m.Source,
			Destination: m.Destination,
			Rw:          m.RW,
		})
	}
	return res
}

func mapNetworks(nets map[string]*network.EndpointSettings) map[string]*protos.NetworkSettings {
	res := make(map[string]*protos.NetworkSettings)
	for name, n := range nets {
		if n == nil {
			continue
		}
		res[name] = &protos.NetworkSettings{
			IpAddress: n.IPAddress,
			Gateway:   n.Gateway,
			NetworkId: n.NetworkID,
		}
	}
	return res
}

func mapStats(stats container.StatsResponse) *protos.GetContainerStatsResponse {
	var rxTotal, txTotal uint64
	for _, netStats := range stats.Networks {
		rxTotal += netStats.RxBytes
		txTotal += netStats.TxBytes
	}

	return &protos.GetContainerStatsResponse{
		Id:        stats.ID,
		CpuUsage:  stats.CPUStats.CPUUsage.TotalUsage,
		MemUsage:  stats.MemoryStats.Usage,
		MemLimit:  stats.MemoryStats.Limit,
		NetInput:  rxTotal,
		NetOutput: txTotal,
	}
}

func mapConfig(req *protos.CreateContainerRequest) *container.Config {
	return &container.Config{
		Image: req.GetImage(),
		Cmd:   req.GetCmd(),
		Env:   req.GetEnv(),
	}
}

func mapHostConfig(h *protos.HostConfig) *container.HostConfig {
	if h == nil {
		return &container.HostConfig{}
	}

	portBindings := nat.PortMap{}
	for _, binding := range h.GetPortBindings() {
		for _, m := range binding.GetHostPorts() {
			port := nat.Port(fmt.Sprintf("%d/tcp", m.ContainerPort))
			hostBinding := nat.PortBinding{
				HostIP:   m.GetHostIp(),
				HostPort: strconv.Itoa(int(m.GetHostPort())),
			}
			portBindings[port] = append(portBindings[port], hostBinding)
		}
	}

	var mounts []mount.Mount
	for _, m := range h.GetMounts() {
		mounts = append(mounts, mount.Mount{
			Source:   m.GetSource(),
			Target:   m.GetTarget(),
			ReadOnly: m.GetReadOnly(),
			Type:     mount.TypeBind,
		})
	}

	return &container.HostConfig{
		PortBindings: portBindings,
		AutoRemove:   h.GetAutoRemove(),
		RestartPolicy: container.RestartPolicy{
			Name:              container.RestartPolicyMode(h.GetRestartPolicy().GetName()),
			MaximumRetryCount: int(h.GetRestartPolicy().GetMaximumRetryCount()),
		},
		Mounts: mounts,
	}
}

func mapNetworking(networks []string) *network.NetworkingConfig {
	if len(networks) == 0 {
		return nil
	}

	endpoints := make(map[string]*network.EndpointSettings, len(networks))
	for _, n := range networks {
		endpoints[n] = &network.EndpointSettings{}
	}

	return &network.NetworkingConfig{EndpointsConfig: endpoints}
}

func mapPorts(ports []container.Port) []string {
	result := make([]string, 0, len(ports))
	for _, port := range ports {
		portStr := fmt.Sprintf("%d/%s", port.PrivatePort, port.Type)
		if port.PublicPort > 0 {
			if port.IP != "" {
				portStr = fmt.Sprintf("%s:%d->%s", port.IP, port.PublicPort, portStr)
			} else {
				portStr = fmt.Sprintf("%d->%s", port.PublicPort, portStr)
			}
		}
		result = append(result, portStr)
	}
	return result
}
