package image

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type mockLayerAPI struct {
	mock.Mock
}

func (m *mockLayerAPI) GetImages(ctx context.Context, opts image.ListOptions) ([]image.Summary, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]image.Summary), args.Error(1)
}

func (m *mockLayerAPI) GetImageDetails(ctx context.Context, id string) (image.InspectResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(image.InspectResponse), args.Error(1)
}

func (m *mockLayerAPI) PullImage(
	ctx context.Context,
	imageName string,
	opts image.PullOptions,
) (io.ReadCloser, error) {
	args := m.Called(ctx, imageName, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *mockLayerAPI) RemoveImage(
	ctx context.Context,
	imageID string,
	opts image.RemoveOptions,
) ([]image.DeleteResponse, error) {
	args := m.Called(ctx, imageID, opts)
	return args.Get(0).([]image.DeleteResponse), args.Error(1)
}

func (m *mockLayerAPI) PruneImage(ctx context.Context, args filters.Args) (image.PruneReport, error) {
	mockArgs := m.Called(ctx, args)
	return mockArgs.Get(0).(image.PruneReport), mockArgs.Error(1)
}

func (m *mockLayerAPI) BuildImage(
	ctx context.Context,
	buildContext io.Reader,
	opts build.ImageBuildOptions,
) (build.ImageBuildResponse, error) {
	args := m.Called(ctx, buildContext, opts)
	return args.Get(0).(build.ImageBuildResponse), args.Error(1)
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

func TestImageServiceRegister(t *testing.T) {
	s := grpc.NewServer()
	NewImageService(nil).Register(s)
	if _, ok := s.GetServiceInfo()["image.v1.ImageService"]; !ok {
		keys := make([]string, 0, len(s.GetServiceInfo()))
		for k := range s.GetServiceInfo() {
			keys = append(keys, k)
		}
		t.Fatalf("expected image.v1.ImageService, registered: %v", keys)
	}
}

func TestServiceGetImages(t *testing.T) {
	tests := []struct {
		name        string
		req         *protos.GetImagesRequest
		setup       func(*mockLayerAPI)
		expectErr   bool
		code        codes.Code
		expectCount int
	}{
		{
			name: "successful images listing",
			req:  &protos.GetImagesRequest{All: true},
			setup: func(ml *mockLayerAPI) {
				images := []image.Summary{
					{ID: "img-1", RepoTags: []string{"nginx:latest"}},
					{ID: "img-2", RepoTags: []string{"redis:alpine"}},
				}
				ml.On("GetImages",
					mock.Anything,
					mock.Anything,
				).Return(images, nil)
			},
			expectErr:   false,
			code:        codes.OK,
			expectCount: 2,
		},
		{
			name: "layer error",
			req:  &protos.GetImagesRequest{All: false},
			setup: func(ml *mockLayerAPI) {
				ml.On("GetImages",
					mock.Anything,
					mock.Anything,
				).Return([]image.Summary{}, errors.New("docker daemon error"))
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
			resp, err := svc.GetImages(context.Background(), tt.req)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.code, grpcCode(err))
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Len(t, resp.Images, tt.expectCount)
			}
			ml.AssertExpectations(t)
		})
	}
}

func TestServiceGetImageDetails(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.GetImageDetailsRequest
		setup     func(*mockLayerAPI)
		expectErr bool
		code      codes.Code
		expectID  string
	}{
		{
			name: "successful image details",
			req:  &protos.GetImageDetailsRequest{Id: "img-123"},
			setup: func(ml *mockLayerAPI) {
				details := image.InspectResponse{
					ID:       "img-123",
					RepoTags: []string{"nginx:latest"},
					Size:     123456,
				}
				ml.On("GetImageDetails",
					mock.Anything,
					"img-123",
				).Return(details, nil)
			},
			expectErr: false,
			code:      codes.OK,
			expectID:  "img-123",
		},
		{
			name:      "missing id error",
			req:       &protos.GetImageDetailsRequest{Id: ""},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name: "layer error",
			req:  &protos.GetImageDetailsRequest{Id: "invalid-id"},
			setup: func(ml *mockLayerAPI) {
				ml.On("GetImageDetails",
					mock.Anything,
					"invalid-id",
				).Return(image.InspectResponse{}, errors.New("image not found"))
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
			resp, err := svc.GetImageDetails(context.Background(), tt.req)

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

func TestServiceRemoveImage(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.RemoveImageRequest
		setup     func(*mockLayerAPI)
		expectErr bool
		code      codes.Code
	}{
		{
			name: "successful image removal",
			req:  &protos.RemoveImageRequest{Id: "img-123"},
			setup: func(ml *mockLayerAPI) {
				deleteResp := []image.DeleteResponse{
					{Deleted: "img-123"},
				}
				ml.On("RemoveImage",
					mock.Anything,
					"img-123",
					mock.Anything,
				).Return(deleteResp, nil)
			},
			expectErr: false,
			code:      codes.OK,
		},
		{
			name:      "missing id error",
			req:       &protos.RemoveImageRequest{Id: ""},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name: "layer error",
			req:  &protos.RemoveImageRequest{Id: "invalid-id"},
			setup: func(ml *mockLayerAPI) {
				ml.On("RemoveImage",
					mock.Anything,
					"invalid-id",
					mock.Anything,
				).Return([]image.DeleteResponse{}, errors.New("image not found"))
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
			resp, err := svc.RemoveImage(context.Background(), tt.req)

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

func TestServicePruneImages(t *testing.T) {
	tests := []struct {
		name      string
		req       *protos.PruneImagesRequest
		setup     func(*mockLayerAPI)
		expectErr bool
		code      codes.Code
	}{
		{
			name: "successful images pruning",
			req:  &protos.PruneImagesRequest{},
			setup: func(ml *mockLayerAPI) {
				pruneReport := image.PruneReport{
					ImagesDeleted:  []image.DeleteResponse{{Deleted: "img-1"}},
					SpaceReclaimed: 123456,
				}
				ml.On("PruneImage",
					mock.Anything,
					mock.Anything,
				).Return(pruneReport, nil)
			},
			expectErr: false,
			code:      codes.OK,
		},
		{
			name: "layer error",
			req:  &protos.PruneImagesRequest{},
			setup: func(ml *mockLayerAPI) {
				ml.On("PruneImage",
					mock.Anything,
					mock.Anything,
				).Return(image.PruneReport{}, errors.New("docker daemon error"))
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
			resp, err := svc.PruneImages(context.Background(), tt.req)

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

func TestServicePullImage(t *testing.T) {
	tests := []struct {
		name        string
		req         *protos.PullImageRequest
		setup       func(*mockLayerAPI)
		setupStream func(*mockPullImageStream)
		expectErr   bool
		code        codes.Code
	}{
		{
			name:      "missing image name",
			req:       &protos.PullImageRequest{Link: ""},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name: "successful image pull",
			req:  &protos.PullImageRequest{Link: "nginx:latest"},
			setup: func(ml *mockLayerAPI) {
				mockReader := &mockReadCloser{data: []byte(`{"status":"Pulling from library/nginx"}`)}
				ml.On("PullImage",
					mock.Anything,
					"nginx:latest",
					mock.Anything,
				).Return(mockReader, nil)
			},
			setupStream: func(ms *mockPullImageStream) {
				ms.On("Send", mock.Anything).Return(nil)
			},
			expectErr: false,
			code:      codes.OK,
		},
		{
			name: "layer error",
			req:  &protos.PullImageRequest{Link: "nonexistent:latest"},
			setup: func(ml *mockLayerAPI) {
				ml.On("PullImage",
					mock.Anything,
					"nonexistent:latest",
					mock.Anything,
				).Return(nil, errors.New("image not found"))
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

			mockStream := &mockPullImageStream{}
			if tt.setupStream != nil {
				tt.setupStream(mockStream)
			}

			svc := &Service{layer: ml}

			err := svc.PullImage(tt.req, mockStream)

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

func TestServiceBuildImage(t *testing.T) {
	tests := []struct {
		name        string
		req         *protos.BuildImageRequest
		setup       func(*mockLayerAPI)
		setupStream func(*mockBuildImageStream)
		expectErr   bool
		code        codes.Code
	}{
		{
			name:      "missing dockerfile",
			req:       &protos.BuildImageRequest{Dockerfile: ""},
			expectErr: true,
			code:      codes.InvalidArgument,
		},
		{
			name: "successful image build",
			req: &protos.BuildImageRequest{
				Dockerfile: "FROM alpine:latest\nRUN echo 'hello'",
				Tags:       []string{"test:latest"},
			},
			setup: func(ml *mockLayerAPI) {
				mockBuildResp := build.ImageBuildResponse{
					Body: &mockReadCloser{data: []byte(`{"stream":"Step 1/2 : FROM alpine:latest"}`)},
				}
				ml.On("BuildImage",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(mockBuildResp, nil)
			},
			setupStream: func(ms *mockBuildImageStream) {
				ms.On("Send", mock.Anything).Return(nil)
			},
			expectErr: false,
			code:      codes.OK,
		},
		{
			name: "layer error",
			req: &protos.BuildImageRequest{
				Dockerfile: "FROM nonexistent:latest",
				Tags:       []string{"test:latest"},
			},
			setup: func(ml *mockLayerAPI) {
				ml.On("BuildImage",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(build.ImageBuildResponse{}, errors.New("base image not found"))
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

			mockStream := &mockBuildImageStream{}
			if tt.setupStream != nil {
				tt.setupStream(mockStream)
			}

			svc := &Service{layer: ml}

			err := svc.BuildImage(tt.req, mockStream)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.code, grpcCode(err))
			} else {
				assert.NoError(t, err)
			}

			if tt.setup != nil {
				ml.AssertExpectations(t)
			}
			if tt.setupStream != nil {
				mockStream.AssertExpectations(t)
			}
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

type mockPullImageStream struct {
	mock.Mock
}

func (m *mockPullImageStream) Send(resp *protos.PullImageResponse) error {
	args := m.Called(resp)
	return args.Error(0)
}

func (m *mockPullImageStream) Context() context.Context {
	return context.Background()
}

func (m *mockPullImageStream) SendMsg(msg interface{}) error {
	return nil
}

func (m *mockPullImageStream) RecvMsg(msg interface{}) error {
	return nil
}

func (m *mockPullImageStream) SetHeader(metadata.MD) error {
	return nil
}

func (m *mockPullImageStream) SendHeader(metadata.MD) error {
	return nil
}

func (m *mockPullImageStream) SetTrailer(metadata.MD) {
}

type mockBuildImageStream struct {
	mock.Mock
}

func (m *mockBuildImageStream) Send(resp *protos.BuildImageResponse) error {
	args := m.Called(resp)
	return args.Error(0)
}

func (m *mockBuildImageStream) Context() context.Context {
	return context.Background()
}

func (m *mockBuildImageStream) SendMsg(msg interface{}) error {
	return nil
}

func (m *mockBuildImageStream) RecvMsg(msg interface{}) error {
	return nil
}

func (m *mockBuildImageStream) SetHeader(metadata.MD) error {
	return nil
}

func (m *mockBuildImageStream) SendHeader(metadata.MD) error {
	return nil
}

func (m *mockBuildImageStream) SetTrailer(metadata.MD) {
}
