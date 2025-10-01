// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

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
