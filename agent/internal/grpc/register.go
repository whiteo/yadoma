// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package grpcserver

import "google.golang.org/grpc"

type Registrator interface {
	Register(rpc *grpc.Server)
}

func RegisterAll(rpc *grpc.Server, services ...Registrator) {
	for _, s := range services {
		s.Register(rpc)
	}
}
