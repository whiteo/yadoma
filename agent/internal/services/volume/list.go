// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package volume provides the agent's service layer for Docker volume management.
// It implements gRPC-facing handlers that validate requests, delegate to the
// Docker client layer, map results to protobuf messages, and translate errors
// into gRPC status codes.
//
// Supported operations include creating volumes, listing and inspecting details,
// removing volumes, and pruning unused volumes. Calls respect the caller's
// context and deadlines; streaming endpoints are not used.
//
// The package does not spawn goroutines on behalf of the caller and relies on
// context cancellation for shutdown. It is intended for internal use by the
// agent's gRPC server layer.
package volume

import (
	"context"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/volume"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetVolumes lists Docker volumes using default volume.ListOptions.
// It honors the incoming context for deadlines and cancellation.
// On success, it maps Docker volume fields (name, driver, mountpoint, labels) into a protos.GetVolumesResponse
// and includes any warnings returned by the Docker API.
// On failure, it returns a gRPC error with code Internal.
func (s *Service) GetVolumes(ctx context.Context, _ *protos.GetVolumesRequest) (*protos.GetVolumesResponse, error) {
	vols, err := s.layer.GetVolumes(ctx, volume.ListOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot list volumes: %v", err)
	}

	resp := &protos.GetVolumesResponse{
		Volumes:  make([]*protos.Volume, 0, len(vols.Volumes)),
		Warnings: vols.Warnings,
	}

	for _, v := range vols.Volumes {
		resp.Volumes = append(resp.Volumes, &protos.Volume{
			Name:       v.Name,
			Driver:     v.Driver,
			Mountpoint: v.Mountpoint,
			Labels:     v.Labels,
		})
	}

	return resp, nil
}
