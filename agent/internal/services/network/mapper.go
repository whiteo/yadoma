// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package network provides the agent's service layer for Docker network management.
// It exposes gRPC-facing handlers that validate requests, delegate to the Docker
// layer, map results to protobuf messages, and translate errors into gRPC status
// codes.
//
// Supported operations include creating networks, listing and inspecting details,
// connecting and disconnecting containers, removing networks, and pruning unused
// networks. Calls respect the caller's context and deadlines; streaming endpoints
// are not used.
//
// The package does not spawn goroutines on behalf of the caller and relies on
// context cancellation for shutdown. It is intended for internal use by the
// agent's gRPC server layer.
package network

import (
	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/network"
)

func mapCreateOptions(req *protos.CreateNetworkRequest) network.CreateOptions {
	ipamConfig := network.IPAMConfig{
		Subnet:  req.GetSubnet(),
		Gateway: req.GetGateway(),
	}

	return network.CreateOptions{
		Driver:     req.GetDriver(),
		Internal:   req.GetInternal(),
		Attachable: req.GetAttachable(),
		Ingress:    req.GetIngress(),
		Labels:     req.GetLabels(),
		IPAM: &network.IPAM{
			Driver: "default",
			Config: []network.IPAMConfig{ipamConfig},
		},
	}
}

func mapEndpointSettings(cfg *protos.EndpointSettings) *network.EndpointSettings {
	return &network.EndpointSettings{
		IPAddress:  cfg.GetIpAddress(),
		MacAddress: cfg.GetMacAddress(),
		Aliases:    cfg.GetAliases(),
	}
}
