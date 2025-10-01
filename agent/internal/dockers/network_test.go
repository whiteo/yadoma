// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package docker

import (
	"context"
	"errors"
	"testing"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (m *MockDockerClient) NetworkList(ctx context.Context,
	options network.ListOptions) ([]network.Summary, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]network.Summary), args.Error(1)
}

func (m *MockDockerClient) NetworkInspect(ctx context.Context,
	networkID string,
	options network.InspectOptions) (network.Inspect, error) {
	args := m.Called(ctx, networkID, options)
	return args.Get(0).(network.Inspect), args.Error(1)
}

func (m *MockDockerClient) NetworkCreate(ctx context.Context,
	name string,
	options network.CreateOptions,
) (network.CreateResponse, error) {
	args := m.Called(ctx, name, options)
	return args.Get(0).(network.CreateResponse), args.Error(1)
}

func (m *MockDockerClient) NetworkConnect(ctx context.Context,
	networkID,
	containerID string,
	config *network.EndpointSettings,
) error {
	args := m.Called(ctx, networkID, containerID, config)
	return args.Error(0)
}

func (m *MockDockerClient) NetworkDisconnect(ctx context.Context,
	networkID,
	containerID string,
	force bool,
) error {
	args := m.Called(ctx, networkID, containerID, force)
	return args.Error(0)
}

func (m *MockDockerClient) NetworkRemove(ctx context.Context, networkID string) error {
	args := m.Called(ctx, networkID)
	return args.Error(0)
}

func (m *MockDockerClient) NetworksPrune(ctx context.Context,
	pruneFilters filters.Args,
) (network.PruneReport, error) {
	args := m.Called(ctx, pruneFilters)
	return args.Get(0).(network.PruneReport), args.Error(1)
}

func TestNetworkGetList(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockDockerClient)
		expected    []network.Summary
		expectError bool
	}{
		{
			name: "successful retrieval of network list",
			setupMock: func(m *MockDockerClient) {
				expected := []network.Summary{
					{ID: "network1", Name: "bridge", Driver: "bridge"},
					{ID: "network2", Name: "host", Driver: "host"},
				}
				m.On("NetworkList",
					mock.Anything,
					network.ListOptions{},
				).Return(expected, nil)
			},
			expected: []network.Summary{
				{ID: "network1", Name: "bridge", Driver: "bridge"},
				{ID: "network2", Name: "host", Driver: "host"},
			},
			expectError: false,
		},
		{
			name: "error when receiving network list",
			setupMock: func(m *MockDockerClient) {
				m.On("NetworkList",
					mock.Anything,
					network.ListOptions{},
				).Return([]network.Summary{}, errors.New("docker error"))
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
			result, err := l.GetNetworks(ctx, network.ListOptions{})

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot list network")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestNetworkGetDetails(t *testing.T) {
	tests := []struct {
		name        string
		networkID   string
		setupMock   func(*MockDockerClient)
		expected    network.Inspect
		expectError bool
	}{
		{
			name:      "successful network details retrieval",
			networkID: "test-network-id",
			setupMock: func(m *MockDockerClient) {
				expected := network.Inspect{
					ID:     "test-network-id",
					Name:   "test-network",
					Driver: "bridge",
				}
				m.On("NetworkInspect",
					mock.Anything,
					"test-network-id",
					network.InspectOptions{},
				).Return(expected, nil)
			},
			expected: network.Inspect{
				ID:     "test-network-id",
				Name:   "test-network",
				Driver: "bridge",
			},
			expectError: false,
		},
		{
			name:      "error when getting network details",
			networkID: "invalid-id",
			setupMock: func(m *MockDockerClient) {
				m.On("NetworkInspect",
					mock.Anything,
					"invalid-id",
					network.InspectOptions{},
				).Return(network.Inspect{}, errors.New("network not found"))
			},
			expected:    network.Inspect{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockDockerClient{}
			tt.setupMock(mockClient)

			l := &Layer{client: mockClient}

			ctx := context.Background()
			result, err := l.GetNetworkDetails(ctx, tt.networkID, network.InspectOptions{})

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot inspect network")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestNetworkCreate(t *testing.T) {
	tests := []struct {
		name          string
		networkName   string
		createOptions network.CreateOptions
		setupMock     func(*MockDockerClient)
		expectedID    string
		expectError   bool
	}{
		{
			name:        "successful network creation",
			networkName: "test-network",
			createOptions: network.CreateOptions{
				Driver: "bridge",
			},
			setupMock: func(m *MockDockerClient) {
				expected := network.CreateResponse{
					ID: "new-network-id",
				}
				m.On("NetworkCreate",
					mock.Anything,
					"test-network",
					network.CreateOptions{Driver: "bridge"},
				).Return(expected, nil)
			},
			expectedID:  "new-network-id",
			expectError: false,
		},
		{
			name:        "error when creating network",
			networkName: "invalid-network",
			createOptions: network.CreateOptions{
				Driver: "invalid-driver",
			},
			setupMock: func(m *MockDockerClient) {
				m.On("NetworkCreate",
					mock.Anything,
					"invalid-network",
					network.CreateOptions{Driver: "invalid-driver"},
				).Return(network.CreateResponse{}, errors.New("invalid driver"))
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
			resp, err := l.CreateNetwork(ctx, tt.networkName, tt.createOptions)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot create network")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, resp.ID)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestNetworkRemove(t *testing.T) {
	tests := []struct {
		name        string
		networkID   string
		setupMock   func(*MockDockerClient)
		expectError bool
	}{
		{
			name:      "successful network removal",
			networkID: "test-network-id",
			setupMock: func(m *MockDockerClient) {
				m.On("NetworkRemove",
					mock.Anything,
					"test-network-id",
				).Return(nil)
			},
			expectError: false,
		},
		{
			name:      "error when removing network",
			networkID: "network-in-use",
			setupMock: func(m *MockDockerClient) {
				m.On("NetworkRemove",
					mock.Anything,
					"network-in-use",
				).Return(errors.New("network in use"))
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
			err := l.RemoveNetwork(ctx, tt.networkID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot remove network")
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestNetworkConnect(t *testing.T) {
	tests := []struct {
		name             string
		networkID        string
		containerID      string
		endpointSettings *network.EndpointSettings
		setupMock        func(*MockDockerClient)
		expectError      bool
	}{
		{
			name:        "successful network connection",
			networkID:   "test-network-id",
			containerID: "test-container-id",
			endpointSettings: &network.EndpointSettings{
				IPAMConfig: &network.EndpointIPAMConfig{
					IPv4Address: "172.20.0.10",
				},
			},
			setupMock: func(m *MockDockerClient) {
				m.On("NetworkConnect",
					mock.Anything,
					"test-network-id",
					"test-container-id",
					mock.Anything,
				).Return(nil)
			},
			expectError: false,
		},
		{
			name:             "error when connecting to network",
			networkID:        "invalid-network",
			containerID:      "test-container-id",
			endpointSettings: nil,
			setupMock: func(m *MockDockerClient) {
				m.On("NetworkConnect",
					mock.Anything,
					"invalid-network",
					"test-container-id",
					(*network.EndpointSettings)(nil),
				).Return(errors.New("network not found"))
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
			err := l.ConnectNetwork(ctx, tt.networkID, tt.containerID, tt.endpointSettings)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot connect container")
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestNetworkConnectContextCancelled(t *testing.T) {
	mockClient := &MockDockerClient{}

	mockClient.On("NetworkConnect",
		mock.Anything,
		"n1",
		"c1",
		(*network.EndpointSettings)(nil),
	).Return(context.Canceled)

	l := &Layer{client: mockClient}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := l.ConnectNetwork(ctx, "n1", "c1", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot connect container")

	mockClient.AssertExpectations(t)
}

func TestNetworkDisconnect(t *testing.T) {
	tests := []struct {
		name        string
		networkID   string
		containerID string
		force       bool
		setupMock   func(*MockDockerClient)
		expectError bool
	}{
		{
			name:        "successful network disconnection",
			networkID:   "test-network-id",
			containerID: "test-container-id",
			force:       false,
			setupMock: func(m *MockDockerClient) {
				m.On("NetworkDisconnect",
					mock.Anything,
					"test-network-id",
					"test-container-id",
					false,
				).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "successful forced network disconnection",
			networkID:   "test-network-id",
			containerID: "test-container-id",
			force:       true,
			setupMock: func(m *MockDockerClient) {
				m.On("NetworkDisconnect",
					mock.Anything,
					"test-network-id",
					"test-container-id",
					true,
				).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "error when disconnecting from network",
			networkID:   "invalid-network",
			containerID: "test-container-id",
			force:       false,
			setupMock: func(m *MockDockerClient) {
				m.On("NetworkDisconnect",
					mock.Anything,
					"invalid-network",
					"test-container-id",
					false,
				).Return(errors.New("network not found"))
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
			err := l.DisconnectNetwork(ctx, tt.networkID, tt.containerID, tt.force)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot disconnect container")
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestNetworksPrune(t *testing.T) {
	tests := []struct {
		name         string
		pruneFilters filters.Args
		setupMock    func(*MockDockerClient)
		expectError  bool
	}{
		{
			name:         "successful networks pruning",
			pruneFilters: filters.Args{},
			setupMock: func(m *MockDockerClient) {
				expected := network.PruneReport{
					NetworksDeleted: []string{"network1", "network2"},
				}
				m.On("NetworksPrune",
					mock.Anything,
					filters.Args{},
				).Return(expected, nil)
			},
			expectError: false,
		},
		{
			name:         "error when pruning networks",
			pruneFilters: filters.Args{},
			setupMock: func(m *MockDockerClient) {
				m.On("NetworksPrune",
					mock.Anything,
					filters.Args{},
				).Return(network.PruneReport{}, errors.New("prune error"))
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
			_, err := l.PruneNetworks(ctx, tt.pruneFilters)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot prune networks")
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestNetworkContextTimeout(t *testing.T) {
	mockClient := &MockDockerClient{}

	mockClient.On("NetworkList",
		mock.Anything,
		network.ListOptions{},
	).Return([]network.Summary{}, context.DeadlineExceeded)

	l := &Layer{client: mockClient}

	ctx := context.Background()
	_, err := l.GetNetworks(ctx, network.ListOptions{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot list network")

	mockClient.AssertExpectations(t)
}
