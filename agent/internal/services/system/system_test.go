package system

import (
	"context"
	"errors"
	"testing"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/system"
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

func (m *mockLayerAPI) GetSystemInfo(ctx context.Context) (system.Info, error) {
	args := m.Called(ctx)
	return args.Get(0).(system.Info), args.Error(1)
}

func (m *mockLayerAPI) GetDiskUsage(ctx context.Context, opts types.DiskUsageOptions) (types.DiskUsage, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(types.DiskUsage), args.Error(1)
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

func TestSystemServiceRegister(t *testing.T) {
	s := grpc.NewServer()
	NewSystemService(nil).Register(s)
	if _, ok := s.GetServiceInfo()["system.v1.SystemService"]; !ok {
		keys := make([]string, 0, len(s.GetServiceInfo()))
		for k := range s.GetServiceInfo() {
			keys = append(keys, k)
		}
		t.Fatalf("expected system.v1.SystemService, registered: %v", keys)
	}
}

func TestServiceGetSystemInfo(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.GetSystemInfoRequest
		setup     func(*mockLayerAPI)
		expectErr bool
		code      codes.Code
	}{
		{
			name: "successful system info retrieval",
			req:  &protos.GetSystemInfoRequest{},
			setup: func(ml *mockLayerAPI) {
				info := system.Info{
					ID:                "system-123",
					Containers:        5,
					ContainersRunning: 2,
					ContainersPaused:  1,
					ContainersStopped: 2,
					Images:            10,
					Driver:            "overlay2",
					MemTotal:          8589934592,
					NCPU:              4,
					Name:              "docker-host",
					ServerVersion:     "20.10.17",
				}
				ml.On("GetSystemInfo", mock.Anything).Return(info, nil)
			},
			expectErr: false,
			code:      codes.OK,
		},
		{
			name: "layer error",
			req:  &protos.GetSystemInfoRequest{},
			setup: func(ml *mockLayerAPI) {
				ml.On("GetSystemInfo",
					mock.Anything,
				).Return(system.Info{}, errors.New("docker daemon error"))
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
			resp, err := svc.GetSystemInfo(context.Background(), tt.req)

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

func TestServiceGetDiskUsage(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.GetDiskUsageRequest
		setup     func(*mockLayerAPI)
		expectErr bool
		code      codes.Code
	}{
		{
			name: "successful disk usage retrieval",
			req:  &protos.GetDiskUsageRequest{},
			setup: func(ml *mockLayerAPI) {
				diskUsage := types.DiskUsage{
					LayersSize: 1073741824,
					Images: []*image.Summary{
						{ID: "img-1", Size: 536870912},
						{ID: "img-2", Size: 268435456},
					},
					Containers: []*container.Summary{
						{ID: "cont-1", SizeRw: 1048576},
					},
					Volumes: []*volume.Volume{
						{Name: "vol-1", UsageData: &volume.UsageData{Size: 104857600}},
					},
				}
				ml.On("GetDiskUsage",
					mock.Anything,
					mock.Anything,
				).Return(diskUsage, nil)
			},
			expectErr: false,
			code:      codes.OK,
		},
		{
			name: "layer error",
			req:  &protos.GetDiskUsageRequest{},
			setup: func(ml *mockLayerAPI) {
				ml.On("GetDiskUsage",
					mock.Anything,
					mock.Anything,
				).Return(types.DiskUsage{}, errors.New("docker daemon error"))
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
			resp, err := svc.GetDiskUsage(context.Background(), tt.req)

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
