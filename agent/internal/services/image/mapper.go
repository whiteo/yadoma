// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package image provides the agent's service layer for Docker image management.
// It implements gRPC-facing handlers that validate requests, delegate to the
// Docker layer, map results to protobuf messages, and translate errors into gRPC
// status codes.
//
// Supported operations include building images from a context, pulling from
// registries, listing and inspecting details, removing images, and pruning
// unused images. Streaming endpoints (for example, build and pull progress)
// propagate the caller's context; callers must consume and close returned
// streams.
//
// The package does not spawn goroutines on behalf of the caller and relies on
// context deadlines and cancellation for shutdown. It is intended for internal
// use by the agent's gRPC server layer.
package image

import (
	"bytes"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/build"
)

func mapBuildOptions(req *protos.BuildImageRequest) (build.ImageBuildOptions, *bytes.Reader) {
	buildArgs := make(map[string]*string, len(req.GetBuildArgs()))
	for k, v := range req.GetBuildArgs() {
		val := v
		buildArgs[k] = &val
	}

	opts := build.ImageBuildOptions{
		Tags:        req.GetTags(),
		NoCache:     req.GetNoCache(),
		Dockerfile:  req.GetDockerfile(),
		BuildArgs:   buildArgs,
		Labels:      req.GetLabels(),
		Remove:      true,
		ForceRemove: true,
	}

	return opts, bytes.NewReader(req.GetBuildContext())
}
