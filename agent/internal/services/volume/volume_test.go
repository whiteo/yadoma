package volume

import (
	"context"
	"errors"
	"testing"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockLayerAPI struct {
	mock.Mock
}

func (m *mockLayerAPI) GetVolumes(ctx context.Context, opts volume.ListOptions) (volume.ListResponse, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(volume.ListResponse), args.Error(1)
}

func (m *mockLayerAPI) GetVolumeDetails(ctx context.Context, volumeID string) (volume.Volume, error) {
	args := m.Called(ctx, volumeID)
	return args.Get(0).(volume.Volume), args.Error(1)
}

func (m *mockLayerAPI) CreateVolume(ctx context.Context, opts volume.CreateOptions) (volume.Volume, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(volume.Volume), args.Error(1)
}

func (m *mockLayerAPI) RemoveVolume(ctx context.Context, volumeID string, force bool) error {
	args := m.Called(ctx, volumeID, force)
	return args.Error(0)
}

func (m *mockLayerAPI) PruneVolumes(ctx context.Context, pruneFilters filters.Args) (volume.PruneReport, error) {
	args := m.Called(ctx, pruneFilters)
	return args.Get(0).(volume.PruneReport), args.Error(1)
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

func TestVolumeServiceRegister(t *testing.T) {
	s := grpc.NewServer()
	NewVolumeService(nil).Register(s)
	if _, ok := s.GetServiceInfo()["volume.v1.VolumeService"]; !ok {
		keys := make([]string, 0, len(s.GetServiceInfo()))
		for k := range s.GetServiceInfo() {
			keys = append(keys, k)
		}
		t.Fatalf("expected volume.v1.VolumeService, registered: %v", keys)
	}
}

func TestServiceGetVolumes(t *testing.T) {
	tests := []struct {
		name        string
		req         *protos.GetVolumesRequest
		setup       func(*mockLayerAPI)
		expectErr   bool
		code        codes.Code
		expectCount int
	}{
		{
			name: "successful volumes listing",
			req:  &protos.GetVolumesRequest{},
			setup: func(ml *mockLayerAPI) {
				volumes := volume.ListResponse{
					Volumes: []*volume.Volume{
						{Name: "vol-1", Driver: "local", Mountpoint: "/var/lib/docker/volumes/vol-1/_data"},
						{Name: "vol-2", Driver: "local", Mountpoint: "/var/lib/docker/volumes/vol-2/_data"},
					},
				}
				ml.On("GetVolumes",
					mock.Anything,
					volume.ListOptions{},
				).Return(volumes, nil)
			},
			expectErr:   false,
			code:        codes.OK,
			expectCount: 2,
		},
		{
			name: "layer error",
			req:  &protos.GetVolumesRequest{},
			setup: func(ml *mockLayerAPI) {
				ml.On("GetVolumes",
					mock.Anything,
					volume.ListOptions{},
				).Return(volume.ListResponse{}, errors.New("docker daemon error"))
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
			resp, err := svc.GetVolumes(context.Background(), tt.req)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.code, grpcCode(err))
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Len(t, resp.Volumes, tt.expectCount)
			}
			ml.AssertExpectations(t)
		})
	}
}

func TestServiceGetVolumeDetails(t *testing.T) {
	tests := []struct {
		name       string
		req        *protos.GetVolumeDetailsRequest
		setup      func(*mockLayerAPI)
		expectErr  bool
		code       codes.Code
		expectName string
	}{
		{
			name: "successful volume details",
			req:  &protos.GetVolumeDetailsRequest{Id: "vol-123"},
			setup: func(ml *mockLayerAPI) {
				vol := volume.Volume{
					Name:       "vol-123",
					Driver:     "local",
					Mountpoint: "/var/lib/docker/volumes/vol-123/_data",
					CreatedAt:  "2023-01-01T00:00:00Z",
					Labels:     map[string]string{"env": "test"},
				}
				ml.On("GetVolumeDetails",
					mock.Anything,
					"vol-123",
				).Return(vol, nil)
			},
			expectErr:  false,
			code:       codes.OK,
			expectName: "vol-123",
		},
		{
			name:      "missing id error",
			req:       &protos.GetVolumeDetailsRequest{Id: ""},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name: "layer error",
			req:  &protos.GetVolumeDetailsRequest{Id: "invalid-id"},
			setup: func(ml *mockLayerAPI) {
				ml.On("GetVolumeDetails",
					mock.Anything,
					"invalid-id",
				).Return(volume.Volume{}, errors.New("volume not found"))
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
			resp, err := svc.GetVolumeDetails(context.Background(), tt.req)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.code, grpcCode(err))
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectName, resp.Name)
			}
			ml.AssertExpectations(t)
		})
	}
}

func TestServiceCreateVolume(t *testing.T) {
	tests := []struct {
		name       string
		req        *protos.CreateVolumeRequest
		setup      func(*mockLayerAPI)
		expectErr  bool
		code       codes.Code
		expectName string
	}{
		{
			name: "successful volume creation",
			req: &protos.CreateVolumeRequest{
				Name:   "test-volume",
				Driver: "local",
				Labels: map[string]string{"env": "test"},
			},
			setup: func(ml *mockLayerAPI) {
				vol := volume.Volume{
					Name:       "test-volume",
					Driver:     "local",
					Mountpoint: "/var/lib/docker/volumes/test-volume/_data",
					Labels:     map[string]string{"env": "test"},
				}
				ml.On("CreateVolume",
					mock.Anything,
					mock.Anything,
				).Return(vol, nil)
			},
			expectErr:  false,
			code:       codes.OK,
			expectName: "test-volume",
		},
		{
			name:      "missing name error",
			req:       &protos.CreateVolumeRequest{Name: "", Driver: "local"},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name: "layer error",
			req: &protos.CreateVolumeRequest{
				Name:   "test-volume",
				Driver: "local",
			},
			setup: func(ml *mockLayerAPI) {
				ml.On("CreateVolume",
					mock.Anything,
					mock.Anything,
				).Return(volume.Volume{}, errors.New("volume creation failed"))
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
			resp, err := svc.CreateVolume(context.Background(), tt.req)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.code, grpcCode(err))
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectName, resp.Name)
			}
			ml.AssertExpectations(t)
		})
	}
}

func TestServiceRemoveVolume(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.RemoveVolumeRequest
		setup     func(*mockLayerAPI)
		expectErr bool
		code      codes.Code
	}{
		{
			name: "successful volume removal",
			req:  &protos.RemoveVolumeRequest{Id: "vol-123", Force: false},
			setup: func(ml *mockLayerAPI) {
				ml.On("RemoveVolume",
					mock.Anything,
					"vol-123",
					false,
				).Return(nil)
			},
			expectErr: false,
			code:      codes.OK,
		},
		{
			name: "successful force volume removal",
			req:  &protos.RemoveVolumeRequest{Id: "vol-123", Force: true},
			setup: func(ml *mockLayerAPI) {
				ml.On("RemoveVolume",
					mock.Anything,
					"vol-123",
					true,
				).Return(nil)
			},
			expectErr: false,
			code:      codes.OK,
		},
		{
			name:      "missing id error",
			req:       &protos.RemoveVolumeRequest{Id: "", Force: false},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name: "layer error",
			req:  &protos.RemoveVolumeRequest{Id: "invalid-id", Force: false},
			setup: func(ml *mockLayerAPI) {
				ml.On("RemoveVolume",
					mock.Anything,
					"invalid-id",
					false,
				).Return(errors.New("volume not found"))
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
			resp, err := svc.RemoveVolume(context.Background(), tt.req)

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

func TestServicePruneVolumes(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.PruneVolumesRequest
		setup     func(*mockLayerAPI)
		expectErr bool
		code      codes.Code
	}{
		{
			name: "successful volumes pruning",
			req:  &protos.PruneVolumesRequest{All: "true"},
			setup: func(ml *mockLayerAPI) {
				pruneReport := volume.PruneReport{
					VolumesDeleted: []string{"vol-1", "vol-2"},
					SpaceReclaimed: 1073741824,
				}
				expectedFilters := filters.NewArgs(filters.Arg("All", "true"))
				ml.On("PruneVolumes",
					mock.Anything,
					expectedFilters,
				).Return(pruneReport, nil)
			},
			expectErr: false,
			code:      codes.OK,
		},
		{
			name: "layer error",
			req:  &protos.PruneVolumesRequest{All: "false"},
			setup: func(ml *mockLayerAPI) {
				expectedFilters := filters.NewArgs(filters.Arg("All", "false"))
				ml.On("PruneVolumes",
					mock.Anything,
					expectedFilters,
				).Return(volume.PruneReport{}, errors.New("docker daemon error"))
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
			resp, err := svc.PruneVolumes(context.Background(), tt.req)

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
