package network

import (
	"context"
	"errors"
	"testing"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockLayerAPI struct {
	mock.Mock
}

func (m *mockLayerAPI) GetNetworks(ctx context.Context, opts network.ListOptions) ([]network.Summary, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]network.Summary), args.Error(1)
}

func (m *mockLayerAPI) GetNetworkDetails(ctx context.Context, id string, opts network.InspectOptions) (network.Inspect, error) {
	args := m.Called(ctx, id, opts)
	return args.Get(0).(network.Inspect), args.Error(1)
}

func (m *mockLayerAPI) CreateNetwork(ctx context.Context, name string, opts network.CreateOptions) (network.CreateResponse, error) {
	args := m.Called(ctx, name, opts)
	return args.Get(0).(network.CreateResponse), args.Error(1)
}

func (m *mockLayerAPI) RemoveNetwork(ctx context.Context, networkID string) error {
	args := m.Called(ctx, networkID)
	return args.Error(0)
}

func (m *mockLayerAPI) ConnectNetwork(ctx context.Context, networkID, containerID string, config *network.EndpointSettings) error {
	args := m.Called(ctx, networkID, containerID, config)
	return args.Error(0)
}

func (m *mockLayerAPI) DisconnectNetwork(ctx context.Context, networkID, containerID string, force bool) error {
	args := m.Called(ctx, networkID, containerID, force)
	return args.Error(0)
}

func (m *mockLayerAPI) PruneNetworks(ctx context.Context, args filters.Args) (network.PruneReport, error) {
	mockArgs := m.Called(ctx, args)
	return mockArgs.Get(0).(network.PruneReport), mockArgs.Error(1)
}

func grpcCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	st, ok := status.FromError(err)
	if !ok {
		return codes.Unknown
	}
	return st.Code()
}

func TestNetworkServiceRegister(t *testing.T) {
	s := grpc.NewServer()
	NewNetworkService(nil).Register(s)
	if _, ok := s.GetServiceInfo()["network.v1.NetworkService"]; !ok {
		keys := make([]string, 0, len(s.GetServiceInfo()))
		for k := range s.GetServiceInfo() {
			keys = append(keys, k)
		}
		t.Fatalf("expected network.v1.NetworkService, registered: %v", keys)
	}
}

func TestServiceGetNetworks(t *testing.T) {
	tests := []struct {
		name        string
		req         *protos.GetNetworksRequest
		setup       func(*mockLayerAPI)
		expectErr   bool
		code        codes.Code
		expectCount int
	}{
		{
			name: "successful networks listing",
			req:  &protos.GetNetworksRequest{},
			setup: func(ml *mockLayerAPI) {
				networks := []network.Summary{
					{ID: "net-1", Name: "bridge", Driver: "bridge"},
					{ID: "net-2", Name: "host", Driver: "host"},
				}
				ml.On("GetNetworks",
					mock.Anything,
					mock.Anything,
				).Return(networks, nil)
			},
			expectErr:   false,
			code:        codes.OK,
			expectCount: 2,
		},
		{
			name: "layer error",
			req:  &protos.GetNetworksRequest{},
			setup: func(ml *mockLayerAPI) {
				ml.On("GetNetworks",
					mock.Anything,
					mock.Anything,
				).Return([]network.Summary{}, errors.New("docker daemon error"))
			},
			expectErr: true,
			code:      codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &mockLayerAPI{}
			if tt.setup != nil {
				tt.setup(ml)
			}
			svc := &Service{layer: ml}
			resp, err := svc.GetNetworks(context.Background(), tt.req)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.code, grpcCode(err))
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Len(t, resp.Networks, tt.expectCount)
			}
			ml.AssertExpectations(t)
		})
	}
}

func TestServiceGetNetworkDetails(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.GetNetworkDetailsRequest
		setup     func(*mockLayerAPI)
		expectErr bool
		code      codes.Code
		expectID  string
	}{
		{
			name: "successful network details",
			req:  &protos.GetNetworkDetailsRequest{Id: "net-123"},
			setup: func(ml *mockLayerAPI) {
				details := network.Inspect{
					ID:     "net-123",
					Name:   "test-network",
					Driver: "bridge",
				}
				ml.On("GetNetworkDetails",
					mock.Anything,
					"net-123",
					mock.Anything,
				).Return(details, nil)
			},
			expectErr: false,
			code:      codes.OK,
			expectID:  "net-123",
		},
		{
			name:      "missing id error",
			req:       &protos.GetNetworkDetailsRequest{Id: ""},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name: "layer error",
			req:  &protos.GetNetworkDetailsRequest{Id: "invalid-id"},
			setup: func(ml *mockLayerAPI) {
				ml.On("GetNetworkDetails",
					mock.Anything,
					"invalid-id",
					mock.Anything,
				).Return(network.Inspect{}, errors.New("network not found"))
			},
			expectErr: true,
			code:      codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &mockLayerAPI{}
			if tt.setup != nil {
				tt.setup(ml)
			}
			svc := &Service{layer: ml}
			resp, err := svc.GetNetworkDetails(context.Background(), tt.req)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.code, grpcCode(err))
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectID, resp.Id)
			}
			ml.AssertExpectations(t)
		})
	}
}

func TestServiceCreateNetwork(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.CreateNetworkRequest
		setup     func(*mockLayerAPI)
		expectErr bool
		code      codes.Code
		expectID  string
	}{
		{
			name: "successful network creation",
			req: &protos.CreateNetworkRequest{
				Name:   "test-network",
				Driver: "bridge",
			},
			setup: func(ml *mockLayerAPI) {
				createResp := network.CreateResponse{
					ID:      "net-123",
					Warning: "",
				}
				ml.On("CreateNetwork",
					mock.Anything,
					"test-network",
					mock.Anything,
				).Return(createResp, nil)
			},
			expectErr: false,
			code:      codes.OK,
			expectID:  "net-123",
		},
		{
			name:      "missing name error",
			req:       &protos.CreateNetworkRequest{Name: "", Driver: "bridge"},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name: "layer error",
			req: &protos.CreateNetworkRequest{
				Name:   "test-network",
				Driver: "bridge",
			},
			setup: func(ml *mockLayerAPI) {
				ml.On("CreateNetwork",
					mock.Anything,
					"test-network",
					mock.Anything,
				).Return(network.CreateResponse{}, errors.New("network creation failed"))
			},
			expectErr: true,
			code:      codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &mockLayerAPI{}
			if tt.setup != nil {
				tt.setup(ml)
			}
			svc := &Service{layer: ml}
			resp, err := svc.CreateNetwork(context.Background(), tt.req)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.code, grpcCode(err))
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectID, resp.Id)
			}
			ml.AssertExpectations(t)
		})
	}
}

func TestServiceRemoveNetwork(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.RemoveNetworkRequest
		setup     func(*mockLayerAPI)
		expectErr bool
		code      codes.Code
	}{
		{
			name: "successful network removal",
			req:  &protos.RemoveNetworkRequest{Id: "net-123"},
			setup: func(ml *mockLayerAPI) {
				ml.On("RemoveNetwork",
					mock.Anything,
					"net-123",
				).Return(nil)
			},
			expectErr: false,
			code:      codes.OK,
		},
		{
			name:      "missing id error",
			req:       &protos.RemoveNetworkRequest{Id: ""},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name: "layer error",
			req:  &protos.RemoveNetworkRequest{Id: "invalid-id"},
			setup: func(ml *mockLayerAPI) {
				ml.On("RemoveNetwork",
					mock.Anything,
					"invalid-id",
				).Return(errors.New("network not found"))
			},
			expectErr: true,
			code:      codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &mockLayerAPI{}
			if tt.setup != nil {
				tt.setup(ml)
			}
			svc := &Service{layer: ml}
			resp, err := svc.RemoveNetwork(context.Background(), tt.req)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.code, grpcCode(err))
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			ml.AssertExpectations(t)
		})
	}
}

func TestServiceConnectContainer(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.ConnectNetworkRequest
		setup     func(*mockLayerAPI)
		expectErr bool
		code      codes.Code
	}{
		{
			name: "successful container connection",
			req: &protos.ConnectNetworkRequest{
				NetworkId:   "net-123",
				ContainerId: "container-123",
			},
			setup: func(ml *mockLayerAPI) {
				ml.On("ConnectNetwork",
					mock.Anything,
					"net-123",
					"container-123",
					mock.Anything,
				).Return(nil)
			},
			expectErr: false,
			code:      codes.OK,
		},
		{
			name:      "missing network id error",
			req:       &protos.ConnectNetworkRequest{NetworkId: "", ContainerId: "container-123"},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name:      "missing container id error",
			req:       &protos.ConnectNetworkRequest{NetworkId: "net-123", ContainerId: ""},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name: "layer error",
			req: &protos.ConnectNetworkRequest{
				NetworkId:   "net-123",
				ContainerId: "container-123",
			},
			setup: func(ml *mockLayerAPI) {
				ml.On("ConnectNetwork",
					mock.Anything,
					"net-123",
					"container-123",
					mock.Anything,
				).Return(errors.New("connection failed"))
			},
			expectErr: true,
			code:      codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &mockLayerAPI{}
			if tt.setup != nil {
				tt.setup(ml)
			}
			svc := &Service{layer: ml}
			resp, err := svc.ConnectNetwork(context.Background(), tt.req)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.code, grpcCode(err))
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			ml.AssertExpectations(t)
		})
	}
}

func TestServiceDisconnectContainer(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.DisconnectNetworkRequest
		setup     func(*mockLayerAPI)
		expectErr bool
		code      codes.Code
	}{
		{
			name: "successful container disconnection",
			req: &protos.DisconnectNetworkRequest{
				NetworkId:   "net-123",
				ContainerId: "container-123",
				Force:       false,
			},
			setup: func(ml *mockLayerAPI) {
				ml.On("DisconnectNetwork",
					mock.Anything,
					"net-123",
					"container-123",
					false,
				).Return(nil)
			},
			expectErr: false,
			code:      codes.OK,
		},
		{
			name:      "missing network id error",
			req:       &protos.DisconnectNetworkRequest{NetworkId: "", ContainerId: "container-123"},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name:      "missing container id error",
			req:       &protos.DisconnectNetworkRequest{NetworkId: "net-123", ContainerId: ""},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name: "layer error",
			req: &protos.DisconnectNetworkRequest{
				NetworkId:   "net-123",
				ContainerId: "container-123",
				Force:       true,
			},
			setup: func(ml *mockLayerAPI) {
				ml.On("DisconnectNetwork",
					mock.Anything,
					"net-123",
					"container-123",
					true,
				).Return(errors.New("disconnection failed"))
			},
			expectErr: true,
			code:      codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &mockLayerAPI{}
			if tt.setup != nil {
				tt.setup(ml)
			}
			svc := &Service{layer: ml}
			resp, err := svc.DisconnectNetwork(context.Background(), tt.req)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.code, grpcCode(err))
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			ml.AssertExpectations(t)
		})
	}
}

func TestServicePruneNetworks(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.PruneNetworksRequest
		setup     func(*mockLayerAPI)
		expectErr bool
		code      codes.Code
	}{
		{
			name: "successful networks pruning",
			req:  &protos.PruneNetworksRequest{},
			setup: func(ml *mockLayerAPI) {
				pruneReport := network.PruneReport{
					NetworksDeleted: []string{"net-1", "net-2"},
				}
				ml.On("PruneNetworks",
					mock.Anything,
					mock.Anything,
				).Return(pruneReport, nil)
			},
			expectErr: false,
			code:      codes.OK,
		},
		{
			name: "layer error",
			req:  &protos.PruneNetworksRequest{},
			setup: func(ml *mockLayerAPI) {
				ml.On("PruneNetworks",
					mock.Anything,
					mock.Anything,
				).Return(network.PruneReport{}, errors.New("docker daemon error"))
			},
			expectErr: true,
			code:      codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &mockLayerAPI{}
			if tt.setup != nil {
				tt.setup(ml)
			}
			svc := &Service{layer: ml}
			resp, err := svc.PruneNetworks(context.Background(), tt.req)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.code, grpcCode(err))
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			ml.AssertExpectations(t)
		})
	}
}
