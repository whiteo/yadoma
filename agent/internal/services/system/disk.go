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

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetDiskUsage reports aggregate Docker disk usage for images, containers, volumes, and layer data.
// It delegates to the Docker client layer using the incoming context to honor cancellation and deadlines,
// maps the result to a `protos.GetDiskUsageResponse`, and returns it on success.
// On failure, it returns a gRPC error with `codes.Internal` describing the underlying cause.
func (s *Service) GetDiskUsage(
	ctx context.Context,
	_ *protos.GetDiskUsageRequest,
) (*protos.GetDiskUsageResponse, error) {
	usage, err := s.layer.GetDiskUsage(ctx, types.DiskUsageOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get disk usage info: %v", err)
	}

	images := mapDiskUsageImage(usage.Images)
	containers := mapDiskUsageContainer(usage.Containers)
	volumes := mapDiskUsageVolume(usage.Volumes)

	return &protos.GetDiskUsageResponse{
		LayersSize: usage.LayersSize,
		Images:     images,
		Containers: containers,
		Volumes:    volumes,
	}, nil
}
