// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package system

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
		NCpu:              int32(info.NCPU),
		Containers:        int32(info.Containers),
		ContainersRunning: int32(info.ContainersRunning),
		ContainersPaused:  int32(info.ContainersPaused),
		ContainersStopped: int32(info.ContainersStopped),
		Images:            int32(info.Images),
		ServerVersion:     info.ServerVersion,
		OperatingSystem:   info.OperatingSystem,
		Architecture:      info.Architecture,
		MemTotal:          info.MemTotal,
		Driver:            info.Driver,
		Labels:            info.Labels,
	}, nil
}
