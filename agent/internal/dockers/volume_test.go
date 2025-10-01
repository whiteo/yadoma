// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package docker

import (
	"context"
	"errors"
	"testing"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (m *MockDockerClient) VolumeList(ctx context.Context, opts volume.ListOptions) (volume.ListResponse, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(volume.ListResponse), args.Error(1)
}

func (m *MockDockerClient) VolumeInspect(ctx context.Context, volumeID string) (volume.Volume, error) {
	args := m.Called(ctx, volumeID)
	return args.Get(0).(volume.Volume), args.Error(1)
}

func (m *MockDockerClient) VolumeCreate(ctx context.Context, opts volume.CreateOptions) (volume.Volume, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(volume.Volume), args.Error(1)
}

func (m *MockDockerClient) VolumeRemove(ctx context.Context, volumeID string, force bool) error {
	args := m.Called(ctx, volumeID, force)
	return args.Error(0)
}

func (m *MockDockerClient) VolumesPrune(ctx context.Context, pruneFilters filters.Args) (volume.PruneReport, error) {
	args := m.Called(ctx, pruneFilters)
	return args.Get(0).(volume.PruneReport), args.Error(1)
}

func TestVolumeGetList(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockDockerClient)
		expected    volume.ListResponse
		expectError bool
	}{
		{
			name: "successful retrieval of volume list",
			setupMock: func(m *MockDockerClient) {
				expected := volume.ListResponse{
					Volumes: []*volume.Volume{
						{Name: "volume1", Driver: "local"},
						{Name: "volume2", Driver: "local"},
					},
					Warnings: []string{},
				}
				m.On("VolumeList",
					mock.Anything,
					volume.ListOptions{},
				).Return(expected, nil)
			},
			expected: volume.ListResponse{
				Volumes: []*volume.Volume{
					{Name: "volume1", Driver: "local"},
					{Name: "volume2", Driver: "local"},
				},
				Warnings: []string{},
			},
			expectError: false,
		},
		{
			name: "error when receiving volume list",
			setupMock: func(m *MockDockerClient) {
				m.On("VolumeList",
					mock.Anything,
					volume.ListOptions{},
				).Return(volume.ListResponse{}, errors.New("docker error"))
			},
			expected:    volume.ListResponse{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockDockerClient{}
			tt.setupMock(mockClient)

			l := &Layer{client: mockClient}

			ctx := context.Background()
			result, err := l.GetVolumes(ctx, volume.ListOptions{})

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot list volume")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestVolumeGetDetails(t *testing.T) {
	tests := []struct {
		name        string
		volumeID    string
		setupMock   func(*MockDockerClient)
		expected    volume.Volume
		expectError bool
	}{
		{
			name:     "successful volume details retrieval",
			volumeID: "test-volume-id",
			setupMock: func(m *MockDockerClient) {
				expected := volume.Volume{
					Name:       "test-volume",
					Driver:     "local",
					Mountpoint: "/var/lib/docker/volumes/test-volume/_data",
				}
				m.On("VolumeInspect",
					mock.Anything,
					"test-volume-id",
				).Return(expected, nil)
			},
			expected: volume.Volume{
				Name:       "test-volume",
				Driver:     "local",
				Mountpoint: "/var/lib/docker/volumes/test-volume/_data",
			},
			expectError: false,
		},
		{
			name:     "error when getting volume details",
			volumeID: "invalid-id",
			setupMock: func(m *MockDockerClient) {
				m.On("VolumeInspect",
					mock.Anything,
					"invalid-id",
				).Return(volume.Volume{}, errors.New("volume not found"))
			},
			expected:    volume.Volume{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockDockerClient{}
			tt.setupMock(mockClient)

			l := &Layer{client: mockClient}

			ctx := context.Background()
			result, err := l.GetVolumeDetails(ctx, tt.volumeID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot inspect volume")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestVolumeCreate(t *testing.T) {
	tests := []struct {
		name          string
		createOptions volume.CreateOptions
		setupMock     func(*MockDockerClient)
		expected      volume.Volume
		expectError   bool
	}{
		{
			name: "successful volume creation",
			createOptions: volume.CreateOptions{
				Name:   "test-volume",
				Driver: "local",
			},
			setupMock: func(m *MockDockerClient) {
				expected := volume.Volume{
					Name:       "test-volume",
					Driver:     "local",
					Mountpoint: "/var/lib/docker/volumes/test-volume/_data",
				}
				m.On("VolumeCreate",
					mock.Anything,
					volume.CreateOptions{
						Name:   "test-volume",
						Driver: "local",
					},
				).Return(expected, nil)
			},
			expected: volume.Volume{
				Name:       "test-volume",
				Driver:     "local",
				Mountpoint: "/var/lib/docker/volumes/test-volume/_data",
			},
			expectError: false,
		},
		{
			name: "error when creating volume",
			createOptions: volume.CreateOptions{
				Name:   "invalid-volume",
				Driver: "invalid-driver",
			},
			setupMock: func(m *MockDockerClient) {
				m.On("VolumeCreate",
					mock.Anything,
					volume.CreateOptions{
						Name:   "invalid-volume",
						Driver: "invalid-driver",
					},
				).Return(volume.Volume{}, errors.New("invalid driver"))
			},
			expected:    volume.Volume{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockDockerClient{}
			tt.setupMock(mockClient)

			l := &Layer{client: mockClient}

			ctx := context.Background()
			result, err := l.CreateVolume(ctx, tt.createOptions)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot create volume")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestVolumeRemove(t *testing.T) {
	tests := []struct {
		name        string
		volumeID    string
		setupMock   func(*MockDockerClient)
		expectError bool
	}{
		{
			name:     "successful volume removal",
			volumeID: "test-volume-id",
			setupMock: func(m *MockDockerClient) {
				// В volume.go всегда передается force=true
				m.On("VolumeRemove",
					mock.Anything,
					"test-volume-id",
					true,
				).Return(nil)
			},
			expectError: false,
		},
		{
			name:     "error when removing volume",
			volumeID: "volume-in-use",
			setupMock: func(m *MockDockerClient) {
				// В volume.go всегда передается force=true
				m.On("VolumeRemove",
					mock.Anything,
					"volume-in-use",
					true,
				).Return(errors.New("volume in use"))
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
			err := l.RemoveVolume(ctx, tt.volumeID, true)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot remove volume")
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestVolumeRemoveContextCancelled(t *testing.T) {
	mockClient := &MockDockerClient{}

	mockClient.On("VolumeRemove",
		mock.Anything,
		"v1",
		true,
	).Return(context.Canceled)

	l := &Layer{client: mockClient}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := l.RemoveVolume(ctx, "v1", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot remove volume")

	mockClient.AssertExpectations(t)
}

func TestVolumesPrune(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockDockerClient)
		expected    volume.PruneReport
		expectError bool
	}{
		{
			name: "successful volumes pruning",
			setupMock: func(m *MockDockerClient) {
				expected := volume.PruneReport{
					VolumesDeleted: []string{"volume1", "volume2"},
					SpaceReclaimed: 2048,
				}
				// В volume.go всегда используется filters.NewArgs()
				m.On("VolumesPrune",
					mock.Anything,
					filters.NewArgs(),
				).Return(expected, nil)
			},
			expected: volume.PruneReport{
				VolumesDeleted: []string{"volume1", "volume2"},
				SpaceReclaimed: 2048,
			},
			expectError: false,
		},
		{
			name: "error when pruning volumes",
			setupMock: func(m *MockDockerClient) {
				// В volume.go всегда используется filters.NewArgs()
				m.On("VolumesPrune",
					mock.Anything,
					filters.NewArgs(),
				).Return(volume.PruneReport{}, errors.New("prune error"))
			},
			expected:    volume.PruneReport{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockDockerClient{}
			tt.setupMock(mockClient)

			l := &Layer{client: mockClient}

			ctx := context.Background()
			result, err := l.PruneVolumes(ctx, filters.NewArgs())

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot prune volumes")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestVolumeContextTimeout(t *testing.T) {
	mockClient := &MockDockerClient{}

	mockClient.On("VolumeList",
		mock.Anything,
		volume.ListOptions{},
	).Return(volume.ListResponse{}, context.DeadlineExceeded)

	l := &Layer{client: mockClient}

	ctx := context.Background()
	_, err := l.GetVolumes(ctx, volume.ListOptions{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot list volume")

	mockClient.AssertExpectations(t)
}
