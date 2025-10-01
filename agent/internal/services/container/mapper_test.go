// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package container

import (
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/whiteo/yadoma/internal/protos"
)

func TestFormatPorts(t *testing.T) {
	tests := []struct {
		name     string
		ports    []container.Port
		expected []string
	}{
		{
			name:     "no ports",
			ports:    []container.Port{},
			expected: []string{},
		},
		{
			name: "private port only",
			ports: []container.Port{
				{PrivatePort: 8080, Type: "tcp"},
			},
			expected: []string{"8080/tcp"},
		},
		{
			name: "public and private port",
			ports: []container.Port{
				{PrivatePort: 8080, PublicPort: 1221, Type: "tcp"},
			},
			expected: []string{"1221->8080/tcp"},
		},
		{
			name: "public and private port with IP",
			ports: []container.Port{
				{PrivatePort: 8080, PublicPort: 1221, Type: "tcp", IP: "127.0.0.1"},
			},
			expected: []string{"127.0.0.1:1221->8080/tcp"},
		},
		{
			name: "multiple ports",
			ports: []container.Port{
				{PrivatePort: 8080, PublicPort: 1221, Type: "tcp"},
				{PrivatePort: 3000, Type: "tcp"},
				{PrivatePort: 443, PublicPort: 8443, Type: "tcp", IP: "0.0.0.0"},
			},
			expected: []string{"1221->8080/tcp", "3000/tcp", "0.0.0.0:8443->443/tcp"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapPorts(tt.ports)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractStatus(t *testing.T) {
	tests := []struct {
		name     string
		state    *container.State
		expected string
	}{
		{
			name:     "nil state",
			state:    nil,
			expected: "",
		},
		{
			name: "running state",
			state: &container.State{
				Status: "running",
			},
			expected: "running",
		},
		{
			name: "stopped state",
			state: &container.State{
				Status: "exited",
			},
			expected: "exited",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractStatus(tt.state)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapMounts(t *testing.T) {
	tests := []struct {
		name     string
		mounts   []container.MountPoint
		expected []*protos.MountPoint
	}{
		{
			name:     "empty mounts",
			mounts:   []container.MountPoint{},
			expected: []*protos.MountPoint{},
		},
		{
			name: "single mount",
			mounts: []container.MountPoint{
				{
					Source:      "/host/path",
					Destination: "/container/path",
					RW:          true,
				},
			},
			expected: []*protos.MountPoint{
				{
					Source:      "/host/path",
					Destination: "/container/path",
					Rw:          true,
				},
			},
		},
		{
			name: "multiple mounts",
			mounts: []container.MountPoint{
				{
					Source:      "/host/data",
					Destination: "/data",
					RW:          true,
				},
				{
					Source:      "/host/config",
					Destination: "/config",
					RW:          false,
				},
			},
			expected: []*protos.MountPoint{
				{
					Source:      "/host/data",
					Destination: "/data",
					Rw:          true,
				},
				{
					Source:      "/host/config",
					Destination: "/config",
					Rw:          false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapMounts(tt.mounts)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapNetworks(t *testing.T) {
	tests := []struct {
		name     string
		nets     map[string]*network.EndpointSettings
		expected map[string]*protos.NetworkSettings
	}{
		{
			name:     "empty networks",
			nets:     map[string]*network.EndpointSettings{},
			expected: map[string]*protos.NetworkSettings{},
		},
		{
			name: "single network",
			nets: map[string]*network.EndpointSettings{
				"bridge": {
					IPAddress: "172.17.0.2",
					Gateway:   "172.17.0.1",
					NetworkID: "network123",
				},
			},
			expected: map[string]*protos.NetworkSettings{
				"bridge": {
					IpAddress: "172.17.0.2",
					Gateway:   "172.17.0.1",
					NetworkId: "network123",
				},
			},
		},
		{
			name: "multiple networks with nil endpoint",
			nets: map[string]*network.EndpointSettings{
				"bridge": {
					IPAddress: "172.17.0.2",
					Gateway:   "172.17.0.1",
					NetworkID: "network123",
				},
				"custom": nil,
			},
			expected: map[string]*protos.NetworkSettings{
				"bridge": {
					IpAddress: "172.17.0.2",
					Gateway:   "172.17.0.1",
					NetworkId: "network123",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapNetworks(tt.nets)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapStats(t *testing.T) {
	tests := []struct {
		name     string
		stats    container.StatsResponse
		expected *protos.GetContainerStatsResponse
	}{
		{
			name: "basic stats",
			stats: container.StatsResponse{
				ID: "container123",
				CPUStats: container.CPUStats{
					CPUUsage: container.CPUUsage{
						TotalUsage: 1000000000,
					},
				},
				MemoryStats: container.MemoryStats{
					Usage: 128000000,
					Limit: 256000000,
				},
				Networks: map[string]container.NetworkStats{
					"eth0": {
						RxBytes: 1024,
						TxBytes: 2048,
					},
				},
			},
			expected: &protos.GetContainerStatsResponse{
				Id:        "container123",
				CpuUsage:  1000000000,
				MemUsage:  128000000,
				MemLimit:  256000000,
				NetInput:  1024,
				NetOutput: 2048,
			},
		},
		{
			name: "multiple network interfaces",
			stats: container.StatsResponse{
				ID: "container456",
				CPUStats: container.CPUStats{
					CPUUsage: container.CPUUsage{
						TotalUsage: 2000000000,
					},
				},
				MemoryStats: container.MemoryStats{
					Usage: 64000000,
					Limit: 128000000,
				},
				Networks: map[string]container.NetworkStats{
					"eth0": {
						RxBytes: 1024,
						TxBytes: 512,
					},
					"eth1": {
						RxBytes: 2048,
						TxBytes: 1024,
					},
				},
			},
			expected: &protos.GetContainerStatsResponse{
				Id:        "container456",
				CpuUsage:  2000000000,
				MemUsage:  64000000,
				MemLimit:  128000000,
				NetInput:  3072, // 1024 + 2048
				NetOutput: 1536, // 512 + 1024
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapStats(tt.stats)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapConfig(t *testing.T) {
	tests := []struct {
		name     string
		req      *protos.CreateContainerRequest
		expected *container.Config
	}{
		{
			name: "basic config",
			req: &protos.CreateContainerRequest{
				Image: "nginx:latest",
				Cmd:   []string{"nginx", "-g", "daemon off;"},
				Env:   []string{"ENV=prod", "DEBUG=false"},
			},
			expected: &container.Config{
				Image: "nginx:latest",
				Cmd:   []string{"nginx", "-g", "daemon off;"},
				Env:   []string{"ENV=prod", "DEBUG=false"},
			},
		},
		{
			name: "minimal config",
			req: &protos.CreateContainerRequest{
				Image: "alpine:latest",
			},
			expected: &container.Config{
				Image: "alpine:latest",
				Cmd:   nil,
				Env:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapConfig(tt.req)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapHostConfig(t *testing.T) {
	tests := []struct {
		name     string
		req      *protos.HostConfig
		expected *container.HostConfig
	}{
		{
			name:     "nil host config",
			req:      nil,
			expected: &container.HostConfig{},
		},
		{
			name: "host config with port bindings",
			req: &protos.HostConfig{
				PortBindings: map[string]*protos.PortBinding{
					"8080/tcp": {
						HostPorts: []*protos.PortMapping{
							{
								HostIp:        "0.0.0.0",
								HostPort:      8080,
								ContainerPort: 80,
							},
						},
					},
				},
				AutoRemove: true,
				RestartPolicy: &protos.RestartPolicy{
					Name:              "always",
					MaximumRetryCount: 3,
				},
			},
			expected: &container.HostConfig{
				PortBindings: nat.PortMap{
					"80/tcp": []nat.PortBinding{
						{
							HostIP:   "0.0.0.0",
							HostPort: "8080",
						},
					},
				},
				AutoRemove: true,
				RestartPolicy: container.RestartPolicy{
					Name:              "always",
					MaximumRetryCount: 3,
				},
				Mounts: nil,
			},
		},
		{
			name: "host config with mounts and no restart policy",
			req: &protos.HostConfig{
				Mounts: []*protos.Mount{
					{
						Source:   "/host/data",
						Target:   "/data",
						ReadOnly: false,
					},
				},
				RestartPolicy: &protos.RestartPolicy{
					Name:              "",
					MaximumRetryCount: 0,
				},
			},
			expected: &container.HostConfig{
				PortBindings: nat.PortMap{},
				Mounts: []mount.Mount{
					{
						Type:     mount.TypeBind,
						Source:   "/host/data",
						Target:   "/data",
						ReadOnly: false,
					},
				},
				RestartPolicy: container.RestartPolicy{
					Name:              "",
					MaximumRetryCount: 0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapHostConfig(tt.req)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapNetworking(t *testing.T) {
	tests := []struct {
		name     string
		networks []string
		expected *network.NetworkingConfig
	}{
		{
			name:     "empty networks",
			networks: []string{},
			expected: nil,
		},
		{
			name:     "single network",
			networks: []string{"bridge"},
			expected: &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					"bridge": {},
				},
			},
		},
		{
			name:     "multiple networks",
			networks: []string{"bridge", "custom"},
			expected: &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					"bridge": {},
					"custom": {},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapNetworking(tt.networks)
			assert.Equal(t, tt.expected, result)
		})
	}
}
