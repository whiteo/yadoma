//go:build tools

package yadoma

//go:generate protoc --go_out=./agent/internal/protos --go-grpc_out=./agent/internal/protos  ./proto/container.proto
//go:generate protoc --go_out=./agent/internal/protos --go-grpc_out=./agent/internal/protos  ./proto/image.proto
//go:generate protoc --go_out=./agent/internal/protos --go-grpc_out=./agent/internal/protos  ./proto/network.proto
//go:generate protoc --go_out=./agent/internal/protos --go-grpc_out=./agent/internal/protos  ./proto/system.proto
//go:generate protoc --go_out=./agent/internal/protos --go-grpc_out=./agent/internal/protos  ./proto/volume.proto
