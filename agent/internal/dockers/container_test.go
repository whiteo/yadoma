// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package docker

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDockerClient struct {
	mock.Mock
}

const testContainerID = "test-container-id"

func (m *MockDockerClient) ContainerList(ctx context.Context,
	options container.ListOptions,
) ([]container.Summary, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]container.Summary), args.Error(1)
}

func (m *MockDockerClient) ContainerInspect(ctx context.Context,
	containerID string,
) (container.InspectResponse, error) {
	args := m.Called(ctx, containerID)
	return args.Get(0).(container.InspectResponse), args.Error(1)
}

func (m *MockDockerClient) ContainerLogs(ctx context.Context,
	containerID string,
	options container.LogsOptions,
) (io.ReadCloser, error) {
	args := m.Called(ctx, containerID, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockDockerClient) ContainerStats(ctx context.Context,
	containerID string,
	stream bool,
) (container.StatsResponseReader, error) {
	args := m.Called(ctx, containerID, stream)
	return args.Get(0).(container.StatsResponseReader), args.Error(1)
}

func (m *MockDockerClient) ContainerCreate(ctx context.Context,
	config *container.Config,
	hostConfig *container.HostConfig,
	networkingConfig *network.NetworkingConfig,
	platform *ocispec.Platform,
	containerName string,
) (container.CreateResponse, error) {
	args := m.Called(ctx, config, hostConfig, networkingConfig, platform, containerName)
	return args.Get(0).(container.CreateResponse), args.Error(1)
}

func (m *MockDockerClient) ContainerRemove(ctx context.Context,
	containerID string,
	options container.RemoveOptions,
) error {
	args := m.Called(ctx, containerID, options)
	return args.Error(0)
}

func (m *MockDockerClient) ContainerStart(ctx context.Context,
	containerID string,
	options container.StartOptions,
) error {
	args := m.Called(ctx, containerID, options)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (m *MockDockerClient) ContainerStop(ctx context.Context,
	containerID string,
	options container.StopOptions,
) error {
	args := m.Called(ctx, containerID, options)
	err := args.Error(0)
	if err != nil {
		return err
	}
	return nil
}

func (m *MockDockerClient) ContainerRestart(ctx context.Context,
	containerID string,
	options container.StopOptions,
) error {
	args := m.Called(ctx, containerID, options)
	return args.Error(0)
}

func (m *MockDockerClient) ContainerPause(ctx context.Context, containerID string) error {
	args := m.Called(ctx, containerID)
	if err := args.Error(0); err != nil {
		return err
	}
	return nil
}

func (m *MockDockerClient) ContainerUnpause(ctx context.Context, containerID string) error {
	args := m.Called(ctx, containerID)
	result := args.Error(0)
	return result
}

func (m *MockDockerClient) ContainerKill(ctx context.Context, containerID string, signal string) error {
	args := m.Called(ctx, containerID, signal)
	if args.Error(0) == nil {
		return nil
	}
	return args.Error(0)
}

func (m *MockDockerClient) ContainerRename(ctx context.Context, containerID string, newName string) error {
	args := m.Called(ctx, containerID, newName)
	errorResult := args.Error(0)
	return errorResult
}

type MockReadCloser struct {
	*strings.Reader
}

func (m *MockReadCloser) Close() error {
	return nil
}

func TestContainerGetList(t *testing.T) {
	expectedContainerList := []container.Summary{
		{ID: "container1", Names: []string{"/test1"}},
		{ID: "container2", Names: []string{"/test2"}},
	}

	tests := []struct {
		expected    []container.Summary
		name        string
		setupMock   func(*MockDockerClient)
		expectError bool
	}{
		{
			name: "successful retrieval of the container list",
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerList",
					mock.Anything,
					container.ListOptions{All: true},
				).Return(expectedContainerList, nil)
			},
			expected:    expectedContainerList,
			expectError: false,
		},
		{
			name: "error when receiving the list of containers",
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerList",
					mock.Anything, container.ListOptions{All: true},
				).Return([]container.Summary{}, errors.New("docker error"))
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
			result, err := l.GetContainers(ctx, container.ListOptions{All: true})

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot list containers")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestContainerGetDetails(t *testing.T) {
	expectedInspectResponse := container.InspectResponse{
		ContainerJSONBase: &container.ContainerJSONBase{
			ID:   testContainerID,
			Name: "/test-container",
		},
	}

	tests := []struct {
		setupMock   func(*MockDockerClient)
		expected    container.InspectResponse
		name        string
		containerID string
		expectError bool
	}{
		{
			name:        "successful container details retrieval",
			containerID: testContainerID,
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerInspect",
					mock.Anything,
					testContainerID,
				).Return(expectedInspectResponse, nil)
			},
			expected:    expectedInspectResponse,
			expectError: false,
		},
		{
			name:        "error when getting container details",
			containerID: "invalid-id",
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerInspect",
					mock.Anything,
					"invalid-id",
				).Return(container.InspectResponse{}, errors.New("container not found"))
			},
			expected:    container.InspectResponse{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockDockerClient{}
			tt.setupMock(mockClient)

			l := &Layer{client: mockClient}

			ctx := context.Background()
			result, err := l.GetContainerDetails(ctx, tt.containerID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot inspect container")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestContainerGetLogs(t *testing.T) {
	expectedOptions := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     true,
		Details:    true,
	}

	tests := []struct {
		setupMock   func(*MockDockerClient)
		name        string
		containerID string
		expectError bool
	}{
		{
			name:        "successful retrieval of container logs",
			containerID: testContainerID,
			setupMock: func(m *MockDockerClient) {
				mockReader := &MockReadCloser{strings.NewReader("test logs")}
				m.On("ContainerLogs",
					mock.Anything,
					testContainerID,
					expectedOptions,
				).Return(io.ReadCloser(mockReader), nil)
			},
			expectError: false,
		},
		{
			name:        "error when receiving container logs",
			containerID: "invalid-id",
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerLogs",
					mock.Anything,
					"invalid-id",
					expectedOptions,
				).Return(io.ReadCloser((*MockReadCloser)(nil)), errors.New("container not found"))
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
			result, err := l.GetContainerLogs(ctx, tt.containerID, expectedOptions)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot get logs for container")
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				err = result.Close()
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestContainerLogsStreamMultiChunksAndMidError(t *testing.T) {
	mockClient := &MockDockerClient{}

	pr, pw := io.Pipe()
	go func() {
		_, _ = pw.Write([]byte("chunk-1\n"))
		_ = pw.CloseWithError(errors.New("stream error"))
	}()

	opts := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: false,
		Follow:     true,
		Details:    false,
	}

	mockClient.On("ContainerLogs",
		mock.Anything,
		"c1",
		opts,
	).Return(io.ReadCloser(pr), nil)

	l := &Layer{client: mockClient}
	ctx := context.Background()

	rc, err := l.GetContainerLogs(ctx, "c1", opts)
	assert.NoError(t, err)

	buf := make([]byte, 16)
	_, rErr1 := rc.Read(buf)
	assert.NoError(t, rErr1)

	_, rErr2 := rc.Read(buf)
	assert.Error(t, rErr2)

	_ = rc.Close()
	mockClient.AssertExpectations(t)
}

func TestContainerLogsContextCancelled(t *testing.T) {
	mockClient := &MockDockerClient{}

	opts := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     true,
		Details:    true,
	}

	mockClient.On("ContainerLogs",
		mock.Anything,
		"c1",
		opts,
	).Return(io.ReadCloser((*MockReadCloser)(nil)), context.Canceled)

	l := &Layer{client: mockClient}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	rc, err := l.GetContainerLogs(ctx, "c1", opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot get logs for container")
	assert.Nil(t, rc)

	mockClient.AssertExpectations(t)
}

func TestContainerGetStats(t *testing.T) {
	tests := []struct {
		setupMock   func(*MockDockerClient)
		name        string
		containerID string
		stream      bool
		expectError bool
	}{
		{
			name:        "successful retrieval of container statistics",
			containerID: testContainerID,
			stream:      true,
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerStats",
					mock.Anything,
					testContainerID,
					true,
				).Return(container.StatsResponseReader{}, nil)
			},
			expectError: false,
		},
		{
			name:        "error when receiving container statistics",
			containerID: "invalid-id",
			stream:      false,
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerStats",
					mock.Anything,
					"invalid-id",
					false,
				).Return(container.StatsResponseReader{}, errors.New("container not found"))
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
			_, err := l.GetContainerStats(ctx, tt.containerID, tt.stream)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot get stats for container")
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestContainerStatsContextCancelled(t *testing.T) {
	mockClient := &MockDockerClient{}

	mockClient.On("ContainerStats",
		mock.Anything,
		"c1",
		true,
	).Return(container.StatsResponseReader{}, context.Canceled)

	l := &Layer{client: mockClient}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := l.GetContainerStats(ctx, "c1", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot get stats for container")

	mockClient.AssertExpectations(t)
}

func TestContainerCreate(t *testing.T) {
	tests := []struct {
		setupMock       func(*MockDockerClient)
		containerConfig *container.Config
		hostConfig      *container.HostConfig
		networkConfig   *network.NetworkingConfig
		platform        *ocispec.Platform
		containerName   string
		name            string
		expectedID      string
		expectError     bool
	}{
		{
			name: "successful container creation",
			containerConfig: &container.Config{
				Image: "nginx:latest",
			},
			hostConfig:    &container.HostConfig{},
			networkConfig: &network.NetworkingConfig{},
			platform:      &ocispec.Platform{},
			containerName: "test-container",
			setupMock: func(m *MockDockerClient) {
				expected := container.CreateResponse{
					ID: "new-container-id",
				}
				m.On("ContainerCreate",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					"test-container",
				).Return(expected, nil)
			},
			expectError: false,
			expectedID:  "new-container-id",
		},
		{
			name: "error when creating a container",
			containerConfig: &container.Config{
				Image: "invalid-image",
			},
			hostConfig:    &container.HostConfig{},
			networkConfig: &network.NetworkingConfig{},
			platform:      &ocispec.Platform{},
			containerName: "test-container",
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerCreate",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					"test-container",
				).Return(container.CreateResponse{}, errors.New("image not found"))
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
			resp, err := l.CreateContainer(ctx, tt.containerConfig, tt.hostConfig, tt.networkConfig, tt.platform, tt.containerName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot create container")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, resp.ID)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestContainerRemove(t *testing.T) {
	tests := []struct {
		setupMock     func(*MockDockerClient)
		name          string
		containerID   string
		removeVolumes bool
		force         bool
		expectError   bool
	}{
		{
			name:          "successful removal of the container",
			containerID:   testContainerID,
			removeVolumes: true,
			force:         false,
			setupMock: func(m *MockDockerClient) {
				expectedOptions := container.RemoveOptions{
					RemoveVolumes: true,
					Force:         false,
				}
				m.On("ContainerRemove",
					mock.Anything,
					testContainerID,
					expectedOptions,
				).Return(nil)
			},
			expectError: false,
		},
		{
			name:          "error when deleting container",
			containerID:   "invalid-id",
			removeVolumes: false,
			force:         true,
			setupMock: func(m *MockDockerClient) {
				expectedOptions := container.RemoveOptions{
					RemoveVolumes: false,
					Force:         true,
				}
				m.On("ContainerRemove",
					mock.Anything,
					"invalid-id",
					expectedOptions,
				).Return(errors.New("container not found"))
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

			opts := container.RemoveOptions{
				RemoveVolumes: tt.removeVolumes,
				Force:         tt.force,
			}

			err := l.RemoveContainer(ctx, tt.containerID, opts)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot remove container")
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestContainerStart(t *testing.T) {
	tests := []struct {
		setupMock   func(*MockDockerClient)
		name        string
		containerID string
		expectError bool
	}{
		{
			name:        "successful container start",
			containerID: testContainerID,
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerStart",
					mock.Anything,
					testContainerID,
					container.StartOptions{},
				).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "error when starting the container",
			containerID: "invalid-id",
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerStart",
					mock.Anything,
					"invalid-id",
					container.StartOptions{},
				).Return(errors.New("container not found"))
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
			err := l.StartContainer(ctx, tt.containerID, container.StartOptions{})

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot start container")
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestContainerStop(t *testing.T) {
	tests := []struct {
		setupMock   func(*MockDockerClient)
		name        string
		containerID string
		expectError bool
	}{
		{
			name:        "successful container stop",
			containerID: testContainerID,
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerStop",
					mock.Anything,
					testContainerID,
					container.StopOptions{},
				).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "error when stopping the container",
			containerID: "invalid-id",
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerStop",
					mock.Anything,
					"invalid-id",
					container.StopOptions{},
				).Return(errors.New("container not found"))
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
			err := l.StopContainer(ctx, tt.containerID, container.StopOptions{})

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot stop container")
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestContainerRestart(t *testing.T) {
	tests := []struct {
		setupMock   func(*MockDockerClient)
		name        string
		containerID string
		expectError bool
	}{
		{
			name:        "successful container restart",
			containerID: testContainerID,
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerRestart",
					mock.Anything,
					testContainerID,
					container.StopOptions{},
				).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "error when restarting the container",
			containerID: "invalid-id",
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerRestart",
					mock.Anything,
					"invalid-id",
					container.StopOptions{},
				).Return(errors.New("container not found"))
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
			err := l.RestartContainer(ctx, tt.containerID, container.StopOptions{})

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot restart container")
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestContainerPause(t *testing.T) {
	tests := []struct {
		setupMock   func(*MockDockerClient)
		name        string
		containerID string
		expectError bool
	}{
		{
			name:        "successful container pause",
			containerID: testContainerID,
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerPause",
					mock.Anything,
					testContainerID,
				).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "error when pausing the container",
			containerID: "invalid-id",
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerPause",
					mock.Anything,
					"invalid-id",
				).Return(errors.New("container not found"))
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
			err := l.PauseContainer(ctx, tt.containerID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot pause container")
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestContainerUnpause(t *testing.T) {
	tests := []struct {
		setupMock   func(*MockDockerClient)
		name        string
		containerID string
		expectError bool
	}{
		{
			name:        "successful container unpause",
			containerID: testContainerID,
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerUnpause",
					mock.Anything,
					testContainerID,
				).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "error when unpausing container",
			containerID: "invalid-id",
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerUnpause",
					mock.Anything,
					"invalid-id",
				).Return(errors.New("container not found"))
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
			err := l.UnpauseContainer(ctx, tt.containerID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot unpause container")
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestContainerKill(t *testing.T) {
	tests := []struct {
		setupMock      func(*MockDockerClient)
		name           string
		containerID    string
		signal         string
		expectedSignal string
		expectError    bool
	}{
		{
			name:           "successful killing of the container with the default signal",
			containerID:    testContainerID,
			signal:         "",
			expectedSignal: "SIGKILL",
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerKill",
					mock.Anything,
					testContainerID,
					"SIGKILL",
				).Return(nil)
			},
			expectError: false,
		},
		{
			name:           "successful killing of container with custom signal",
			containerID:    testContainerID,
			signal:         "SIGTERM",
			expectedSignal: "SIGTERM",
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerKill",
					mock.Anything,
					testContainerID,
					"SIGTERM",
				).Return(nil)
			},
			expectError: false,
		},
		{
			name:           "error when killing the container",
			containerID:    "invalid-id",
			signal:         "SIGKILL",
			expectedSignal: "SIGKILL",
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerKill",
					mock.Anything,
					"invalid-id",
					"SIGKILL",
				).Return(errors.New("container not found"))
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
			err := l.KillContainer(ctx, tt.containerID, tt.signal)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot kill container")
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestContainerRename(t *testing.T) {
	tests := []struct {
		setupMock   func(*MockDockerClient)
		name        string
		containerID string
		newName     string
		expectError bool
	}{
		{
			name:        "successful renaming of the container",
			containerID: testContainerID,
			newName:     "new-container-name",
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerRename",
					mock.Anything,
					testContainerID,
					"new-container-name",
				).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "error when renaming a container",
			containerID: "invalid-id",
			newName:     "new-name",
			setupMock: func(m *MockDockerClient) {
				m.On("ContainerRename",
					mock.Anything,
					"invalid-id",
					"new-name",
				).Return(errors.New("container not found"))
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
			err := l.RenameContainer(ctx, tt.containerID, tt.newName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot rename container")
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestContainerContextTimeout(t *testing.T) {
	mockClient := &MockDockerClient{}

	mockClient.On("ContainerList",
		mock.Anything,
		container.ListOptions{All: true},
	).Return([]container.Summary{}, context.DeadlineExceeded)

	l := &Layer{client: mockClient}

	ctx := context.Background()
	_, err := l.GetContainers(ctx, container.ListOptions{All: true})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot list containers")

	mockClient.AssertExpectations(t)
}
