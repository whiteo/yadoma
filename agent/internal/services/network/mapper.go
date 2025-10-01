// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

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
