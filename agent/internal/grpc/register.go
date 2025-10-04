// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package grpcserver is implemented by services that can register their gRPC handlers on a gRPC server.
package grpcserver

import "google.golang.org/grpc"

type Registrator interface {
	Register(rpc *grpc.Server)
}

// RegisterAll registers each provided service on the given gRPC server.
// It calls s.Register(rpc) for every service, in the order they are passed.
// rpc must be non-nil; passing no services is allowed and results in a no-op.
func RegisterAll(rpc *grpc.Server, services ...Registrator) {
	for _, s := range services {
		s.Register(rpc)
	}
}
