// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package network

import (
	"testing"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/network"
	"github.com/stretchr/testify/assert"
)

func TestMapCreateOptions(t *testing.T) {
	tests := []struct {
		name     string
		req      *protos.CreateNetworkRequest
		expected network.CreateOptions
	}{
		{
			name: "minimal request",
			req: &protos.CreateNetworkRequest{
				Name: "test-network",
			},
			expected: network.CreateOptions{
				Driver:     "",
				Internal:   false,
				Attachable: false,
				Ingress:    false,
				Labels:     nil,
				IPAM: &network.IPAM{
					Driver: "default",
					Config: []network.IPAMConfig{
						{
							Subnet:  "",
							Gateway: "",
						},
					},
				},
			},
		},
		{
			name: "full request",
			req: &protos.CreateNetworkRequest{
				Name:       "prod-network",
				Driver:     "bridge",
				Internal:   true,
				Attachable: true,
				Ingress:    false,
				Subnet:     "172.20.0.0/16",
				Gateway:    "172.20.0.1",
				Labels: map[string]string{
					"env":     "prod",
					"project": "yadoma",
				},
			},
			expected: network.CreateOptions{
				Driver:     "bridge",
				Internal:   true,
				Attachable: true,
				Ingress:    false,
				Labels: map[string]string{
					"env":     "prod",
					"project": "yadoma",
				},
				IPAM: &network.IPAM{
					Driver: "default",
					Config: []network.IPAMConfig{
						{
							Subnet:  "172.20.0.0/16",
							Gateway: "172.20.0.1",
						},
					},
				},
			},
		},
		{
			name: "overlay driver",
			req: &protos.CreateNetworkRequest{
				Name:       "overlay-net",
				Driver:     "overlay",
				Attachable: true,
				Labels:     map[string]string{"type": "overlay"},
			},
			expected: network.CreateOptions{
				Driver:     "overlay",
				Internal:   false,
				Attachable: true,
				Ingress:    false,
				Labels:     map[string]string{"type": "overlay"},
				IPAM: &network.IPAM{
					Driver: "default",
					Config: []network.IPAMConfig{
						{
							Subnet:  "",
							Gateway: "",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapCreateOptions(tt.req)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapEndpointSettings(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *protos.EndpointSettings
		expected *network.EndpointSettings
	}{
		{
			name: "minimal endpoint settings",
			cfg: &protos.EndpointSettings{
				IpAddress: "192.168.1.100",
			},
			expected: &network.EndpointSettings{
				IPAddress:  "192.168.1.100",
				MacAddress: "",
				Aliases:    nil,
			},
		},
		{
			name: "full endpoint settings",
			cfg: &protos.EndpointSettings{
				IpAddress:  "192.168.1.100",
				MacAddress: "02:42:c0:a8:01:64",
				Aliases:    []string{"web", "frontend", "app"},
			},
			expected: &network.EndpointSettings{
				IPAddress:  "192.168.1.100",
				MacAddress: "02:42:c0:a8:01:64",
				Aliases:    []string{"web", "frontend", "app"},
			},
		},
		{
			name: "empty settings",
			cfg: &protos.EndpointSettings{
				IpAddress:  "",
				MacAddress: "",
				Aliases:    []string{},
			},
			expected: &network.EndpointSettings{
				IPAddress:  "",
				MacAddress: "",
				Aliases:    []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapEndpointSettings(tt.cfg)
			assert.Equal(t, tt.expected, result)
		})
	}
}
