// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package docker

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (m *MockDockerClient) ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
	args := m.Called(ctx, options)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]image.Summary), args.Error(1)
}

func (m *MockDockerClient) ImageInspect(ctx context.Context,
	imageID string,
	inspectOpts ...client.ImageInspectOption,
) (image.InspectResponse, error) {
	args := m.Called(ctx, imageID, inspectOpts)
	return args.Get(0).(image.InspectResponse), args.Error(1)
}

func (m *MockDockerClient) ImageRemove(ctx context.Context,
	imageID string,
	options image.RemoveOptions,
) ([]image.DeleteResponse, error) {
	args := m.Called(ctx, imageID, options)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]image.DeleteResponse), args.Error(1)
}

func (m *MockDockerClient) ImagePull(ctx context.Context,
	refStr string,
	options image.PullOptions,
) (io.ReadCloser, error) {
	args := m.Called(ctx, refStr, options)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockDockerClient) ImageBuild(ctx context.Context,
	buildContext io.Reader,
	options build.ImageBuildOptions,
) (build.ImageBuildResponse, error) {
	args := m.Called(ctx, buildContext, options)
	return args.Get(0).(build.ImageBuildResponse), args.Error(1)
}

func (m *MockDockerClient) ImagesPrune(ctx context.Context,
	pruneFilters filters.Args,
) (image.PruneReport, error) {
	args := m.Called(ctx, pruneFilters)
	return args.Get(0).(image.PruneReport), args.Error(1)
}

func TestImageGetList(t *testing.T) {
	tests := []struct {
		expected    []image.Summary
		name        string
		setupMock   func(*MockDockerClient)
		expectError bool
	}{
		{
			name: "successful image list retrieval",
			setupMock: func(m *MockDockerClient) {
				expected := []image.Summary{
					{ID: "image1", RepoTags: []string{"nginx:latest"}},
					{ID: "image2", RepoTags: []string{"redis:alpine"}},
				}
				m.On("ImageList",
					mock.Anything,
					image.ListOptions{All: true},
				).Return(expected, nil)
			},
			expected: []image.Summary{
				{ID: "image1", RepoTags: []string{"nginx:latest"}},
				{ID: "image2", RepoTags: []string{"redis:alpine"}},
			},
			expectError: false,
		},
		{
			name: "error when getting image list",
			setupMock: func(m *MockDockerClient) {
				m.On("ImageList",
					mock.Anything,
					image.ListOptions{All: true},
				).Return(nil, errors.New("docker error"))
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockDockerClient{}
			tt.setupMock(mockClient)

			l := &Layer{client: mockClient}

			ctx := context.Background()
			result, err := l.GetImages(ctx, image.ListOptions{All: true})

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot list image")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestImageGetDetails(t *testing.T) {
	tests := []struct {
		name        string
		imageID     string
		setupMock   func(*MockDockerClient)
		expected    image.InspectResponse
		expectError bool
	}{
		{
			name:    "successful image details retrieval",
			imageID: "test-image-id",
			setupMock: func(m *MockDockerClient) {
				expected := image.InspectResponse{
					ID:       "test-image-id",
					RepoTags: []string{"nginx:latest"},
				}
				m.On("ImageInspect",
					mock.Anything, "test-image-id",
					mock.Anything,
				).Return(expected, nil)
			},
			expected: image.InspectResponse{
				ID:       "test-image-id",
				RepoTags: []string{"nginx:latest"},
			},
			expectError: false,
		},
		{
			name:    "error when getting image details",
			imageID: "invalid-id",
			setupMock: func(m *MockDockerClient) {
				m.On("ImageInspect",
					mock.Anything,
					"invalid-id",
					mock.Anything,
				).Return(image.InspectResponse{}, errors.New("image not found"))
			},
			expected:    image.InspectResponse{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockDockerClient{}
			tt.setupMock(mockClient)

			l := &Layer{client: mockClient}

			ctx := context.Background()
			result, err := l.GetImageDetails(ctx, tt.imageID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot inspect image")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestImageRemove(t *testing.T) {
	tests := []struct {
		name          string
		imageID       string
		setupMock     func(*MockDockerClient)
		expected      []image.DeleteResponse
		pruneChildren bool
		force         bool
		expectError   bool
	}{
		{
			name:          "successful image removal",
			imageID:       "test-image-id",
			pruneChildren: true,
			force:         false,
			setupMock: func(m *MockDockerClient) {
				expected := []image.DeleteResponse{
					{Deleted: "test-image-id"},
				}
				expectedOptions := image.RemoveOptions{
					PruneChildren: true,
					Force:         false,
				}
				m.On("ImageRemove",
					mock.Anything,
					"test-image-id",
					expectedOptions,
				).Return(expected, nil)
			},
			expected: []image.DeleteResponse{
				{Deleted: "test-image-id"},
			},
			expectError: false,
		},
		{
			name:          "error when removing image",
			imageID:       "invalid-id",
			pruneChildren: false,
			force:         true,
			setupMock: func(m *MockDockerClient) {
				expectedOptions := image.RemoveOptions{
					PruneChildren: false,
					Force:         true,
				}
				m.On("ImageRemove",
					mock.Anything,
					"invalid-id",
					expectedOptions,
				).Return(nil, errors.New("image not found"))
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockDockerClient{}
			tt.setupMock(mockClient)

			l := &Layer{client: mockClient}

			ctx := context.Background()

			opts := image.RemoveOptions{
				PruneChildren: tt.pruneChildren,
				Force:         tt.force,
			}

			result, err := l.RemoveImage(ctx, tt.imageID, opts)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot remove image")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestImagePull(t *testing.T) {
	tests := []struct {
		name         string
		refStr       string
		registryAuth string
		setupMock    func(*MockDockerClient)
		expectError  bool
	}{
		{
			name:         "successful image pull",
			refStr:       "nginx:latest",
			registryAuth: "auth-token",
			setupMock: func(m *MockDockerClient) {
				mockReader := &MockReadCloser{strings.NewReader("pull progress")}
				expectedOptions := image.PullOptions{RegistryAuth: "auth-token"}
				m.On("ImagePull",
					mock.Anything,
					"nginx:latest",
					expectedOptions,
				).Return(io.ReadCloser(mockReader), nil)
			},
			expectError: false,
		},
		{
			name:         "error when pulling image",
			refStr:       "invalid-image",
			registryAuth: "",
			setupMock: func(m *MockDockerClient) {
				expectedOptions := image.PullOptions{RegistryAuth: ""}
				m.On("ImagePull",
					mock.Anything,
					"invalid-image",
					expectedOptions,
				).Return(nil, errors.New("image not found"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockDockerClient{}
			tt.setupMock(mockClient)

			l := &Layer{client: mockClient}

			ctx := context.Background()

			opts := image.PullOptions{RegistryAuth: tt.registryAuth}

			result, err := l.PullImage(ctx, tt.refStr, opts)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot pull image")
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				closeErr := result.Close()
				assert.NoError(t, closeErr)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestImagePullStreamMultiChunksAndMidError(t *testing.T) {
	mockClient := &MockDockerClient{}

	pr, pw := io.Pipe()
	go func() {
		_, _ = pw.Write([]byte("part-1\n"))
		_ = pw.CloseWithError(errors.New("stream error"))
	}()

	mockClient.On("ImagePull",
		mock.Anything,
		"nginx:latest",
		image.PullOptions{RegistryAuth: ""},
	).Return(io.NopCloser(pr), nil)

	l := &Layer{client: mockClient}
	ctx := context.Background()

	rc, err := l.PullImage(ctx, "nginx:latest", image.PullOptions{RegistryAuth: ""})
	assert.NoError(t, err)

	buf := make([]byte, 64)
	_, rErr1 := rc.Read(buf)
	assert.NoError(t, rErr1)

	_, rErr2 := rc.Read(buf)
	assert.Error(t, rErr2)

	_ = rc.Close()
	mockClient.AssertExpectations(t)
}

func TestImagePullContextCancelled(t *testing.T) {
	mockClient := &MockDockerClient{}

	mockClient.On("ImagePull",
		mock.Anything,
		"nginx:latest",
		image.PullOptions{RegistryAuth: ""},
	).Return(nil, context.Canceled)

	l := &Layer{client: mockClient}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	rc, err := l.PullImage(ctx, "nginx:latest", image.PullOptions{RegistryAuth: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot pull image")
	assert.Nil(t, rc)

	mockClient.AssertExpectations(t)
}

func TestImageBuild(t *testing.T) {
	tests := []struct {
		name        string
		tag         string
		buildCtx    io.Reader
		setupMock   func(*MockDockerClient)
		expectError bool
	}{
		{
			name:     "successful image build",
			buildCtx: strings.NewReader("dockerfile content"),
			tag:      "my-app:latest",
			setupMock: func(m *MockDockerClient) {
				expected := build.ImageBuildResponse{
					Body: io.NopCloser(strings.NewReader("build output")),
				}
				expectedOptions := build.ImageBuildOptions{
					Tags:        []string{"my-app:latest"},
					Dockerfile:  "Dockerfile",
					Remove:      true,
					ForceRemove: true,
				}
				m.On("ImageBuild",
					mock.Anything,
					mock.Anything,
					expectedOptions,
				).Return(expected, nil)
			},
			expectError: false,
		},
		{
			name:     "error when building image",
			buildCtx: strings.NewReader("invalid dockerfile"),
			tag:      "invalid:tag",
			setupMock: func(m *MockDockerClient) {
				expectedOptions := build.ImageBuildOptions{
					Tags:        []string{"invalid:tag"},
					Dockerfile:  "Dockerfile",
					Remove:      true,
					ForceRemove: true,
				}
				m.On("ImageBuild",
					mock.Anything,
					mock.Anything,
					expectedOptions,
				).Return(build.ImageBuildResponse{}, errors.New("build failed"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockDockerClient{}
			tt.setupMock(mockClient)

			l := &Layer{client: mockClient}

			ctx := context.Background()

			opts := build.ImageBuildOptions{
				Tags:        []string{tt.tag},
				Dockerfile:  "Dockerfile",
				Remove:      true,
				ForceRemove: true,
			}

			result, err := l.BuildImage(ctx, tt.buildCtx, opts)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot build image")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result.Body)
				closeErr := result.Body.Close()
				assert.NoError(t, closeErr)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestImageBuildStreamMultiChunksAndMidError(t *testing.T) {
	mockClient := &MockDockerClient{}

	pr, pw := io.Pipe()
	go func() {
		_, _ = pw.Write([]byte("build-out-1\n"))
		_ = pw.CloseWithError(errors.New("stream error"))
	}()

	opts := build.ImageBuildOptions{
		Tags:        []string{"app:latest"},
		Dockerfile:  "Dockerfile",
		Remove:      true,
		ForceRemove: true,
	}

	mockClient.On("ImageBuild",
		mock.Anything,
		mock.Anything,
		opts,
	).Return(build.ImageBuildResponse{Body: io.NopCloser(pr)}, nil)

	l := &Layer{client: mockClient}
	ctx := context.Background()

	resp, err := l.BuildImage(ctx, strings.NewReader("ctx"), opts)
	assert.NoError(t, err)

	buf := make([]byte, 64)
	_, rErr1 := resp.Body.Read(buf)
	assert.NoError(t, rErr1)

	_, rErr2 := resp.Body.Read(buf)
	assert.Error(t, rErr2)

	_ = resp.Body.Close()
	mockClient.AssertExpectations(t)
}

func TestImageBuildContextCancelled(t *testing.T) {
	mockClient := &MockDockerClient{}

	opts := build.ImageBuildOptions{
		Tags:        []string{"app:latest"},
		Dockerfile:  "Dockerfile",
		Remove:      true,
		ForceRemove: true,
	}

	mockClient.On("ImageBuild",
		mock.Anything,
		mock.Anything,
		opts,
	).Return(build.ImageBuildResponse{}, context.Canceled)

	l := &Layer{client: mockClient}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := l.BuildImage(ctx, strings.NewReader("ctx"), opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot build image")

	mockClient.AssertExpectations(t)
}

func TestImagePrune(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockDockerClient)
		expected    image.PruneReport
		dangling    bool
		expectError bool
	}{
		{
			name:     "successful dangling images prune",
			dangling: true,
			setupMock: func(m *MockDockerClient) {
				expected := image.PruneReport{
					ImagesDeleted:  []image.DeleteResponse{{Deleted: "image1"}},
					SpaceReclaimed: 1024,
				}
				expectedFilters := filters.NewArgs(filters.KeyValuePair{Key: "dangling", Value: "true"})
				m.On("ImagesPrune",
					mock.Anything,
					expectedFilters,
				).Return(expected, nil)
			},
			expected: image.PruneReport{
				ImagesDeleted:  []image.DeleteResponse{{Deleted: "image1"}},
				SpaceReclaimed: 1024,
			},
			expectError: false,
		},
		{
			name:     "successful all unused images prune",
			dangling: false,
			setupMock: func(m *MockDockerClient) {
				expected := image.PruneReport{
					ImagesDeleted:  []image.DeleteResponse{{Deleted: "image2"}},
					SpaceReclaimed: 2048,
				}
				expectedFilters := filters.NewArgs(filters.KeyValuePair{Key: "dangling", Value: "false"})
				m.On("ImagesPrune",
					mock.Anything,
					expectedFilters,
				).Return(expected, nil)
			},
			expected: image.PruneReport{
				ImagesDeleted:  []image.DeleteResponse{{Deleted: "image2"}},
				SpaceReclaimed: 2048,
			},
			expectError: false,
		},
		{
			name:     "error when pruning images",
			dangling: true,
			setupMock: func(m *MockDockerClient) {
				expectedFilters := filters.NewArgs(filters.KeyValuePair{Key: "dangling", Value: "true"})
				m.On("ImagesPrune",
					mock.Anything,
					expectedFilters,
				).Return(image.PruneReport{}, errors.New("prune failed"))
			},
			expected:    image.PruneReport{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockDockerClient{}
			tt.setupMock(mockClient)

			l := &Layer{client: mockClient}

			ctx := context.Background()

			var d string
			if tt.dangling {
				d = "true"
			} else {
				d = "false"
			}

			result, err := l.PruneImage(ctx, filters.NewArgs(filters.KeyValuePair{Key: "dangling", Value: d}))

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot prune image")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestImageContextTimeout(t *testing.T) {
	mockClient := &MockDockerClient{}

	mockClient.On("ImageList",
		mock.Anything,
		image.ListOptions{All: true},
	).Return(nil, context.DeadlineExceeded)

	l := &Layer{client: mockClient}

	ctx := context.Background()
	_, err := l.GetImages(ctx, image.ListOptions{All: true})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot list image")

	mockClient.AssertExpectations(t)
}
