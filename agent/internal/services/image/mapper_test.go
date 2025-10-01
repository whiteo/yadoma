// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package image

import (
	"io"
	"testing"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/build"
	"github.com/stretchr/testify/assert"
)

func TestMapBuildOptions(t *testing.T) {
	tests := []struct {
		name     string
		req      *protos.BuildImageRequest
		validate func(t *testing.T, opts build.ImageBuildOptions, reader io.Reader)
	}{
		{
			name: "minimal request",
			req: &protos.BuildImageRequest{
				Tags:         []string{"test:latest"},
				BuildContext: []byte("FROM alpine"),
			},
			validate: func(t *testing.T, opts build.ImageBuildOptions, reader io.Reader) {
				assert.Equal(t, []string{"test:latest"}, opts.Tags)
				assert.Equal(t, false, opts.NoCache)
				assert.Equal(t, "", opts.Dockerfile)
				assert.Equal(t, true, opts.Remove)
				assert.Equal(t, true, opts.ForceRemove)
				assert.Empty(t, opts.BuildArgs)
				assert.Empty(t, opts.Labels)

				data, err := io.ReadAll(reader)
				assert.NoError(t, err)
				assert.Equal(t, []byte("FROM alpine"), data)
			},
		},
		{
			name: "full request",
			req: &protos.BuildImageRequest{
				Tags:         []string{"test:latest", "test:v1.0"},
				NoCache:      true,
				Dockerfile:   "Dockerfile.prod",
				BuildContext: []byte("FROM alpine\nRUN echo hello"),
				BuildArgs: map[string]string{
					"VERSION": "1.0",
					"ENV":     "prod",
				},
				Labels: map[string]string{
					"version": "1.0",
					"app":     "test",
				},
			},
			validate: func(t *testing.T, opts build.ImageBuildOptions, reader io.Reader) {
				assert.Equal(t, []string{"test:latest", "test:v1.0"}, opts.Tags)
				assert.Equal(t, true, opts.NoCache)
				assert.Equal(t, "Dockerfile.prod", opts.Dockerfile)
				assert.Equal(t, true, opts.Remove)
				assert.Equal(t, true, opts.ForceRemove)

				assert.Len(t, opts.BuildArgs, 2)
				assert.NotNil(t, opts.BuildArgs["VERSION"])
				assert.Equal(t, "1.0", *opts.BuildArgs["VERSION"])
				assert.NotNil(t, opts.BuildArgs["ENV"])
				assert.Equal(t, "prod", *opts.BuildArgs["ENV"])

				assert.Equal(t, map[string]string{
					"version": "1.0",
					"app":     "test",
				}, opts.Labels)

				data, err := io.ReadAll(reader)
				assert.NoError(t, err)
				assert.Equal(t, []byte("FROM alpine\nRUN echo hello"), data)
			},
		},
		{
			name: "empty build args and labels",
			req: &protos.BuildImageRequest{
				Tags:         []string{"test:empty"},
				BuildContext: []byte("FROM scratch"),
				BuildArgs:    map[string]string{},
				Labels:       map[string]string{},
			},
			validate: func(t *testing.T, opts build.ImageBuildOptions, reader io.Reader) {
				assert.Equal(t, []string{"test:empty"}, opts.Tags)
				assert.Len(t, opts.BuildArgs, 0)
				assert.Len(t, opts.Labels, 0)

				data, err := io.ReadAll(reader)
				assert.NoError(t, err)
				assert.Equal(t, []byte("FROM scratch"), data)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildOpts, reader := mapBuildOptions(tt.req)
			tt.validate(t, buildOpts, reader)
		})
	}
}
