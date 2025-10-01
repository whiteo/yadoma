// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package docker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/system"
	"github.com/docker/docker/api/types/volume"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (m *MockDockerClient) Info(ctx context.Context) (system.Info, error) {
	args := m.Called(ctx)
	return args.Get(0).(system.Info), args.Error(1)
}

func (m *MockDockerClient) DiskUsage(ctx context.Context, options types.DiskUsageOptions) (types.DiskUsage, error) {
	args := m.Called(ctx, options)
	return args.Get(0).(types.DiskUsage), args.Error(1)
}

func TestSystemGetInfo(t *testing.T) {
	expectedSystemInfo := system.Info{
		ID:                "test-docker-id",
		Containers:        5,
		ContainersRunning: 3,
		ContainersPaused:  1,
		ContainersStopped: 1,
		Images:            10,
		Driver:            "overlay2",
		MemoryLimit:       true,
		SwapLimit:         true,
		CPUCfsPeriod:      true,
		CPUCfsQuota:       true,
		CPUShares:         true,
		CPUSet:            true,
		PidsLimit:         true,
		IPv4Forwarding:    true,
		Debug:             false,
		NFd:               42,
		NGoroutines:       100,
		SystemTime:        "2023-01-01T00:00:00Z",
		LoggingDriver:     "json-file",
		CgroupDriver:      "cgroupfs",
		NEventsListener:   0,
		KernelVersion:     "5.4.0",
		OperatingSystem:   "Ubuntu 20.04",
		OSType:            "linux",
		Architecture:      "x86_64",
		NCPU:              4,
		MemTotal:          8589934592, // 8GB
		DockerRootDir:     "/var/lib/docker",
		Name:              "test-docker-host",
	}

	tests := []struct {
		name        string
		setupMock   func(*MockDockerClient)
		expected    system.Info
		expectError bool
	}{
		{
			name: "successful system info retrieval",
			setupMock: func(m *MockDockerClient) {
				m.On("Info", mock.Anything).Return(expectedSystemInfo, nil)
			},
			expected:    expectedSystemInfo,
			expectError: false,
		},
		{
			name: "error when getting system info",
			setupMock: func(m *MockDockerClient) {
				m.On("Info",
					mock.Anything,
				).Return(system.Info{}, errors.New("docker daemon not running"))
			},
			expected:    system.Info{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockDockerClient{}
			tt.setupMock(mockClient)

			l := &Layer{client: mockClient}

			ctx := context.Background()
			result, err := l.GetSystemInfo(ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot get system info")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestSystemGetDiskUsage(t *testing.T) {
	lastUsedAt := time.Date(2022, 1, 1, 0, 10, 0, 0, time.UTC)

	expectedDiskUsage := types.DiskUsage{
		LayersSize: 1024000000,
		Images: []*image.Summary{
			{
				ID:      "image1",
				Size:    500000000,
				Created: 1640995200,
			},
			{
				ID:      "image2",
				Size:    524288000,
				Created: 1640995300,
			},
		},
		Containers: []*container.Summary{
			{
				ID:         "container1",
				Names:      []string{"/test1"},
				Image:      "image1",
				SizeRw:     10485760,
				SizeRootFs: 500000000,
			},
		},
		Volumes: []*volume.Volume{
			{
				Name:       "volume1",
				Driver:     "local",
				Mountpoint: "/var/lib/docker/volumes/volume1/_data",
				UsageData: &volume.UsageData{
					Size:     52428800,
					RefCount: 1,
				},
			},
		},
		BuildCache: []*build.CacheRecord{
			{
				ID:         "cache1",
				Type:       "regular",
				Size:       104857600,
				CreatedAt:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				LastUsedAt: &lastUsedAt,
				UsageCount: 5,
				InUse:      false,
				Shared:     false,
			},
		},
	}

	tests := []struct {
		name        string
		setupMock   func(*MockDockerClient)
		expected    types.DiskUsage
		expectError bool
	}{
		{
			name: "successful disk usage retrieval",
			setupMock: func(m *MockDockerClient) {
				m.On("DiskUsage",
					mock.Anything,
					types.DiskUsageOptions{},
				).Return(expectedDiskUsage, nil)
			},
			expected:    expectedDiskUsage,
			expectError: false,
		},
		{
			name: "error when getting disk usage",
			setupMock: func(m *MockDockerClient) {
				m.On("DiskUsage",
					mock.Anything,
					types.DiskUsageOptions{},
				).Return(types.DiskUsage{}, errors.New("access denied"))
			},
			expected:    types.DiskUsage{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockDockerClient{}
			tt.setupMock(mockClient)

			l := &Layer{client: mockClient}

			ctx := context.Background()
			result, err := l.GetDiskUsage(ctx, types.DiskUsageOptions{})

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot get disk usage")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestSystemContextTimeout(t *testing.T) {
	mockClient := &MockDockerClient{}

	mockClient.On("Info", mock.Anything).Return(system.Info{}, context.DeadlineExceeded)

	l := &Layer{client: mockClient}

	ctx := context.Background()
	_, err := l.GetSystemInfo(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot get system info")

	mockClient.AssertExpectations(t)
}
