// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package system provides the agent's service layer for Docker system operations.
// It exposes gRPC-facing handlers that validate requests, delegate to the Docker
// client layer, map results to protobuf messages, and translate errors into gRPC
// status codes.
//
// Supported operations include retrieving daemon/system information and reporting
// aggregate disk usage across images, containers, volumes, and layer sizes.
// Calls respect the caller's context and deadlines; streaming endpoints are not used.
//
// The package does not spawn goroutines on behalf of the caller and relies on
// context cancellation for shutdown. It is intended for internal use by the
// agent's gRPC server layer.
package system

import (
	"context"
	"math"

	"github.com/whiteo/yadoma/internal/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetSystemInfo retrieves Docker daemon and host metadata.
// It respects the incoming context for cancellation/deadlines, delegates to the client layer
// via s.layer.GetSystemInfo, and maps the result to the protobuf response.
// On failure, it returns a gRPC codes.Internal error with additional context.
func (s *Service) GetSystemInfo(
	ctx context.Context,
	_ *protos.GetSystemInfoRequest,
) (*protos.GetSystemInfoResponse, error) {
	info, err := s.layer.GetSystemInfo(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get system info: %v", err)
	}

	return &protos.GetSystemInfoResponse{
		Id:                info.ID,
		Name:              info.Name,
		KernelVersion:     info.KernelVersion,
		NCpu:              clampToInt32(info.NCPU),
		Containers:        clampToInt32(info.Containers),
		ContainersRunning: clampToInt32(info.ContainersRunning),
		ContainersPaused:  clampToInt32(info.ContainersPaused),
		ContainersStopped: clampToInt32(info.ContainersStopped),
		Images:            clampToInt32(info.Images),
		ServerVersion:     info.ServerVersion,
		OperatingSystem:   info.OperatingSystem,
		Architecture:      info.Architecture,
		MemTotal:          info.MemTotal,
		Driver:            info.Driver,
		Labels:            info.Labels,
	}, nil
}

func clampToInt32(v int) int32 {
	if v > math.MaxInt32 {
		return int32(math.MaxInt32)
	}
	if v < math.MinInt32 {
		return int32(math.MinInt32)
	}
	return int32(v)
}
