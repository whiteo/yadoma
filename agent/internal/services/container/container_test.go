package container

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type MockLayer struct{ mock.Mock }

func (m *MockLayer) GetContainers(ctx context.Context, opts container.ListOptions) ([]container.Summary, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]container.Summary), args.Error(1)
}
func (m *MockLayer) GetContainerDetails(ctx context.Context, id string) (container.InspectResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(container.InspectResponse), args.Error(1)
}
func (m *MockLayer) GetContainerLogs(ctx context.Context, id string,
	opts container.LogsOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, id, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}
func (m *MockLayer) GetContainerStats(ctx context.Context,
	id string, stream bool) (container.StatsResponseReader, error) {
	args := m.Called(ctx, id, stream)
	if args.Get(0) == nil {
		return container.StatsResponseReader{}, args.Error(1)
	}
	return args.Get(0).(container.StatsResponseReader), args.Error(1)
}
func (m *MockLayer) CreateContainer(ctx context.Context, c *container.Config,
	h *container.HostConfig, n *network.NetworkingConfig,
	p *ocispec.Platform, name string) (container.CreateResponse, error) {
	args := m.Called(ctx, c, h, n, p, name)
	return args.Get(0).(container.CreateResponse), args.Error(1)
}
func (m *MockLayer) RemoveContainer(ctx context.Context, id string, opts container.RemoveOptions) error {
	args := m.Called(ctx, id, opts)
	return args.Error(0)
}
func (m *MockLayer) StartContainer(ctx context.Context, id string, opts container.StartOptions) error {
	args := m.Called(ctx, id, opts)
	return args.Error(0)
}
func (m *MockLayer) StopContainer(ctx context.Context, id string, opts container.StopOptions) error {
	args := m.Called(ctx, id, opts)
	return args.Error(0)
}
func (m *MockLayer) RestartContainer(ctx context.Context, id string, opts container.StopOptions) error {
	args := m.Called(ctx, id, opts)
	return args.Error(0)
}
func (m *MockLayer) PauseContainer(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockLayer) UnpauseContainer(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockLayer) KillContainer(ctx context.Context, id, signal string) error {
	args := m.Called(ctx, id, signal)
	return args.Error(0)
}
func (m *MockLayer) RenameContainer(ctx context.Context, id, name string) error {
	args := m.Called(ctx, id, name)
	return args.Error(0)
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

func TestContainerServiceRegister(t *testing.T) {
	s := grpc.NewServer()
	NewContainerService(nil).Register(s)
	if _, ok := s.GetServiceInfo()["container.v1.ContainerService"]; !ok {
		keys := make([]string, 0, len(s.GetServiceInfo()))
		for k := range s.GetServiceInfo() {
			keys = append(keys, k)
		}
		t.Fatalf("expected container.v1.ContainerService, registered: %v", keys)
	}
}

func TestServiceGetContainers(t *testing.T) {
	tests := []struct {
		name        string
		req         *protos.GetContainersRequest
		setup       func(*MockLayer)
		expectErr   bool
		code        codes.Code
		expectCount int
	}{
		{
			name: "success list",
			req:  &protos.GetContainersRequest{All: true},
			setup: func(ml *MockLayer) {
				ml.On("GetContainers",
					mock.Anything,
					mock.Anything,
				).Return([]container.Summary{{ID: "c1", Names: []string{"/a"}}}, nil)
			},
			expectErr:   false,
			code:        codes.OK,
			expectCount: 1,
		},
		{
			name: "layer error",
			req:  &protos.GetContainersRequest{All: false},
			setup: func(ml *MockLayer) {
				ml.On("GetContainers",
					mock.Anything,
					mock.Anything,
				).Return([]container.Summary{}, errors.New("boom"))
			},
			expectErr: true,
			code:      codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MockLayer{}
			if tt.setup != nil {
				tt.setup(ml)
			}
			svc := &Service{layer: ml}
			resp, err := svc.GetContainers(context.Background(), tt.req)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.code, grpcCode(err))
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectCount, len(resp.GetContainers()))
			}
			ml.AssertExpectations(t)
		})
	}
}

func TestServiceGetContainerDetails(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.GetContainerDetailsRequest
		setup     func(*MockLayer)
		expectErr bool
		code      codes.Code
		expectID  string
	}{
		{
			name:      "missing id",
			req:       &protos.GetContainerDetailsRequest{Id: ""},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name: "success",
			req:  &protos.GetContainerDetailsRequest{Id: "c1"},
			setup: func(ml *MockLayer) {
				inspect := container.InspectResponse{
					ContainerJSONBase: &container.ContainerJSONBase{ID: "c1", Name: "/c1", Image: "img",
						State: &container.State{Status: "running"}},
					NetworkSettings: &container.NetworkSettings{Networks: map[string]*network.EndpointSettings{}}}
				ml.On("GetContainerDetails", mock.Anything, "c1").Return(inspect, nil)
			},
			expectErr: false,
			code:      codes.OK,
			expectID:  "c1",
		},
		{
			name: "layer error",
			req:  &protos.GetContainerDetailsRequest{Id: "c2"},
			setup: func(ml *MockLayer) {
				ml.On("GetContainerDetails",
					mock.Anything, "c2").Return(container.InspectResponse{}, errors.New("fail"))
			},
			expectErr: true,
			code:      codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MockLayer{}
			if tt.setup != nil {
				tt.setup(ml)
			}
			svc := &Service{layer: ml}
			resp, err := svc.GetContainerDetails(context.Background(), tt.req)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.code, grpcCode(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectID, resp.GetId())
			}
			ml.AssertExpectations(t)
		})
	}
}

func TestServiceCreateContainer(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.CreateContainerRequest
		setup     func(*MockLayer)
		expectErr bool
		code      codes.Code
		expectID  string
	}{
		{
			name: "successful container creation",
			req: &protos.CreateContainerRequest{
				Image: "nginx:latest",
				Name:  "test-container",
			},
			setup: func(ml *MockLayer) {
				ml.On("CreateContainer",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					"test-container",
				).Return(container.CreateResponse{ID: "container-123"}, nil)
			},
			expectErr: false,
			code:      codes.OK,
			expectID:  "container-123",
		},
		{
			name: "missing image error",
			req: &protos.CreateContainerRequest{
				Image: "",
				Name:  "test-container",
			},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name: "layer error",
			req: &protos.CreateContainerRequest{
				Image: "nginx:latest",
				Name:  "test-container",
			},
			setup: func(ml *MockLayer) {
				ml.On("CreateContainer",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					"test-container",
				).Return(container.CreateResponse{}, errors.New("docker error"))
			},
			expectErr: true,
			code:      codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MockLayer{}
			if tt.setup != nil {
				tt.setup(ml)
			}
			svc := &Service{layer: ml}
			resp, err := svc.CreateContainer(context.Background(), tt.req)

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

func TestServiceStateChangingOps(t *testing.T) {
	type op struct {
		name     string
		call     func(s *Service, id string) error
		withName bool
		withOpts interface{}
		layerOn  func(ml *MockLayer, id string)
		missing  bool
	}

	ops := []op{
		{
			name: "start",
			call: func(s *Service, id string) error {
				_, err := s.StartContainer(context.Background(), &protos.StartContainerRequest{Id: id})
				return err
			},
			layerOn: func(ml *MockLayer, id string) {
				ml.On("StartContainer",
					mock.Anything, id, container.StartOptions{}).Return(nil)
			},
		},
		{
			name: "stop",
			call: func(s *Service, id string) error {
				_, err := s.StopContainer(context.Background(), &protos.StopContainerRequest{Id: id})
				return err
			},
			layerOn: func(ml *MockLayer, id string) {
				ml.On("StopContainer",
					mock.Anything, id, container.StopOptions{}).Return(nil)
			},
		},
		{
			name: "restart",
			call: func(s *Service, id string) error {
				_, err := s.RestartContainer(context.Background(), &protos.RestartContainerRequest{Id: id})
				return err
			},
			layerOn: func(ml *MockLayer, id string) {
				ml.On("RestartContainer",
					mock.Anything, id, container.StopOptions{}).Return(nil)
			},
		},
		{
			name: "pause",
			call: func(s *Service, id string) error {
				_, err := s.PauseContainer(context.Background(), &protos.PauseContainerRequest{Id: id})
				return err
			},
			layerOn: func(ml *MockLayer, id string) {
				ml.On("PauseContainer", mock.Anything, id).Return(nil)
			},
		},
		{
			name: "unpause",
			call: func(s *Service, id string) error {
				_, err := s.UnpauseContainer(context.Background(), &protos.UnpauseContainerRequest{Id: id})
				return err
			},
			layerOn: func(ml *MockLayer, id string) {
				ml.On("UnpauseContainer", mock.Anything, id).Return(nil)
			},
		},
		{
			name: "kill",
			call: func(s *Service, id string) error {
				_, err := s.KillContainer(context.Background(), &protos.KillContainerRequest{Id: id})
				return err
			},
			layerOn: func(ml *MockLayer, id string) {
				ml.On("KillContainer", mock.Anything, id, "").Return(nil)
			},
		},
		{
			name: "rename",
			call: func(s *Service, id string) error {
				_, err := s.RenameContainer(context.Background(), &protos.RenameContainerRequest{Id: id, Name: "n"})
				return err
			},
			layerOn: func(ml *MockLayer, id string) {
				ml.On("RenameContainer", mock.Anything, id, "n").Return(nil)
			},
		},
		{
			name: "remove",
			call: func(s *Service, id string) error {
				_, err := s.RemoveContainer(context.Background(), &protos.RemoveContainerRequest{Id: id})
				return err
			},
			layerOn: func(ml *MockLayer, id string) {
				ml.On("RemoveContainer", mock.Anything, id, container.RemoveOptions{}).Return(nil)
			},
		},
	}

	for _, op := range ops {
		t.Run(op.name+" success", func(t *testing.T) {
			ml := &MockLayer{}
			op.layerOn(ml, "c1")
			svc := &Service{layer: ml}
			err := op.call(svc, "c1")
			assert.NoError(t, err)
			ml.AssertExpectations(t)
		})

		t.Run(op.name+" invalid id", func(t *testing.T) {
			ml := &MockLayer{}
			svc := &Service{layer: ml}
			err := op.call(svc, "")
			assert.Equal(t, codes.InvalidArgument, grpcCode(err))
		})

		t.Run(op.name+" layer error", func(t *testing.T) {
			ml := &MockLayer{}
			switch op.name {
			case "start":
				ml.On("StartContainer",
					mock.Anything, "c2", container.StartOptions{}).Return(errors.New("e"))
			case "stop":
				ml.On("StopContainer",
					mock.Anything, "c2", container.StopOptions{}).Return(errors.New("e"))
			case "restart":
				ml.On("RestartContainer",
					mock.Anything, "c2", container.StopOptions{}).Return(errors.New("e"))
			case "pause":
				ml.On("PauseContainer", mock.Anything, "c2").Return(errors.New("e"))
			case "unpause":
				ml.On("UnpauseContainer", mock.Anything, "c2").Return(errors.New("e"))
			case "kill":
				ml.On("KillContainer", mock.Anything, "c2", "").Return(errors.New("e"))
			case "rename":
				ml.On("RenameContainer", mock.Anything, "c2", "n").Return(errors.New("e"))
			case "remove":
				ml.On("RemoveContainer",
					mock.Anything, "c2", container.RemoveOptions{}).Return(errors.New("e"))
			}
			svc := &Service{layer: ml}
			err := op.call(svc, "c2")
			assert.Equal(t, codes.Internal, grpcCode(err))
			ml.AssertExpectations(t)
		})
	}
}

func TestServiceGetContainerLogs(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.GetContainerLogsRequest
		setup     func(*MockLayer, *mockContainerLogsStream)
		expectErr bool
		code      codes.Code
	}{
		{
			name:      "missing id",
			req:       &protos.GetContainerLogsRequest{Id: ""},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name: "successful logs retrieval - non-TTY container",
			req:  &protos.GetContainerLogsRequest{Id: "c1", Follow: true},
			setup: func(ml *MockLayer, stream *mockContainerLogsStream) {
				inspect := container.InspectResponse{
					ContainerJSONBase: &container.ContainerJSONBase{
						ID: "c1",
					},
					Config: &container.Config{
						Tty: false,
					},
				}
				ml.On("GetContainerDetails", mock.Anything, "c1").Return(inspect, nil)

				logLine := "test log line\n"
				header := []byte{
					1, 0, 0, 0,
					0, 0, 0, byte(len(logLine)),
				}
				multiplexedData := append(header, []byte(logLine)...)
				mockReader := &mockReadCloser{data: multiplexedData}

				ml.On("GetContainerLogs",
					mock.Anything,
					"c1",
					mock.MatchedBy(func(opts container.LogsOptions) bool {
						return opts.Follow == true && opts.ShowStdout == true && opts.ShowStderr == true
					}),
				).Return(mockReader, nil)

				stream.On("Send", mock.Anything).Return(nil)
			},
			expectErr: false,
			code:      codes.OK,
		},
		{
			name: "successful logs retrieval - TTY container",
			req:  &protos.GetContainerLogsRequest{Id: "c1-tty", Follow: false},
			setup: func(ml *MockLayer, stream *mockContainerLogsStream) {
				inspect := container.InspectResponse{
					ContainerJSONBase: &container.ContainerJSONBase{
						ID: "c1-tty",
					},
					Config: &container.Config{
						Tty: true,
					},
				}
				ml.On("GetContainerDetails", mock.Anything, "c1-tty").Return(inspect, nil)

				mockReader := &mockReadCloser{data: []byte("raw log output\n")}
				ml.On("GetContainerLogs",
					mock.Anything,
					"c1-tty",
					mock.MatchedBy(func(opts container.LogsOptions) bool {
						return opts.Follow == false && opts.ShowStdout == true && opts.ShowStderr == true
					}),
				).Return(mockReader, nil)

				stream.On("Send", mock.Anything).Return(nil)
			},
			expectErr: false,
			code:      codes.OK,
		},
		{
			name: "layer error - inspect fails",
			req:  &protos.GetContainerLogsRequest{Id: "c2", Follow: false},
			setup: func(ml *MockLayer, stream *mockContainerLogsStream) {
				ml.On("GetContainerDetails", mock.Anything, "c2").Return(
					container.InspectResponse{}, errors.New("container not found"))
			},
			expectErr: true,
			code:      codes.Internal,
		},
		{
			name: "layer error - logs retrieval fails",
			req:  &protos.GetContainerLogsRequest{Id: "c3", Follow: false},
			setup: func(ml *MockLayer, stream *mockContainerLogsStream) {
				inspect := container.InspectResponse{
					ContainerJSONBase: &container.ContainerJSONBase{ID: "c3"},
					Config:            &container.Config{Tty: false},
				}
				ml.On("GetContainerDetails", mock.Anything, "c3").Return(inspect, nil)
				ml.On("GetContainerLogs",
					mock.Anything,
					"c3",
					mock.Anything,
				).Return(nil, errors.New("cannot get logs"))
			},
			expectErr: true,
			code:      codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MockLayer{}
			mockStream := &mockContainerLogsStream{}

			if tt.setup != nil {
				tt.setup(ml, mockStream)
			}

			svc := &Service{layer: ml}
			err := svc.GetContainerLogs(tt.req, mockStream)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.code, grpcCode(err))
			} else {
				assert.NoError(t, err)
			}
			ml.AssertExpectations(t)
			mockStream.AssertExpectations(t)
		})
	}
}

func TestServiceGetContainerStats(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.GetContainerStatsRequest
		setup     func(*MockLayer, *mockContainerStatsStream)
		expectErr bool
		code      codes.Code
	}{
		{
			name:      "missing id",
			req:       &protos.GetContainerStatsRequest{Id: ""},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name: "successful stats retrieval",
			req:  &protos.GetContainerStatsRequest{Id: "c1", Stream: true},
			setup: func(ml *MockLayer, stream *mockContainerStatsStream) {
				mockReader := &mockReadCloser{
					data: []byte(`{"read":"2024-01-01T00:00:00Z","cpu_stats":{"cpu_usage":{"total_usage":1000}}}`),
				}
				statsReader := container.StatsResponseReader{
					Body: mockReader,
				}
				ml.On("GetContainerStats",
					mock.Anything,
					"c1",
					true,
				).Return(statsReader, nil)

				stream.On("Send", mock.Anything).Return(nil)
			},
			expectErr: false,
			code:      codes.OK,
		},
		{
			name: "layer error",
			req:  &protos.GetContainerStatsRequest{Id: "c2", Stream: false},
			setup: func(ml *MockLayer, stream *mockContainerStatsStream) {
				ml.On("GetContainerStats",
					mock.Anything,
					"c2",
					false,
				).Return(container.StatsResponseReader{}, errors.New("container not found"))
			},
			expectErr: true,
			code:      codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := &MockLayer{}
			mockStream := &mockContainerStatsStream{}

			if tt.setup != nil {
				tt.setup(ml, mockStream)
			}

			svc := &Service{layer: ml}
			err := svc.GetContainerStats(tt.req, mockStream)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.code, grpcCode(err))
			} else {
				assert.NoError(t, err)
			}
			ml.AssertExpectations(t)
			mockStream.AssertExpectations(t)
		})
	}
}

type mockReadCloser struct {
	data []byte
	pos  int
}

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	if m.pos >= len(m.data) {
		return 0, io.EOF
	}
	n = copy(p, m.data[m.pos:])
	m.pos += n
	return n, nil
}

func (m *mockReadCloser) Close() error {
	return nil
}

type mockContainerLogsStream struct {
	mock.Mock
}

func (m *mockContainerLogsStream) Send(resp *protos.GetContainerLogsResponse) error {
	args := m.Called(resp)
	return args.Error(0)
}

func (m *mockContainerLogsStream) Context() context.Context {
	return context.Background()
}

func (m *mockContainerLogsStream) SendMsg(msg interface{}) error {
	return nil
}

func (m *mockContainerLogsStream) RecvMsg(msg interface{}) error {
	return nil
}

func (m *mockContainerLogsStream) SetHeader(metadata.MD) error {
	return nil
}

func (m *mockContainerLogsStream) SendHeader(metadata.MD) error {
	return nil
}

func (m *mockContainerLogsStream) SetTrailer(metadata.MD) {
}

type mockContainerStatsStream struct {
	mock.Mock
}

func (m *mockContainerStatsStream) Send(resp *protos.GetContainerStatsResponse) error {
	args := m.Called(resp)
	return args.Error(0)
}

func (m *mockContainerStatsStream) Context() context.Context {
	return context.Background()
}

func (m *mockContainerStatsStream) SendMsg(msg interface{}) error {
	return nil
}

func (m *mockContainerStatsStream) RecvMsg(msg interface{}) error {
	return nil
}

func (m *mockContainerStatsStream) SetHeader(metadata.MD) error {
	return nil
}

func (m *mockContainerStatsStream) SendHeader(metadata.MD) error {
	return nil
}

func (m *mockContainerStatsStream) SetTrailer(metadata.MD) {
}
