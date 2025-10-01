// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package system

import (
	"testing"

	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	"github.com/stretchr/testify/assert"
)

func TestMapDiskUsageImage(t *testing.T) {
	tests := []struct {
		name     string
		imgs     []*image.Summary
		expected []*protos.DiskUsageImage
	}{
		{
			name:     "empty images",
			imgs:     []*image.Summary{},
			expected: []*protos.DiskUsageImage{},
		},
		{
			name: "single image",
			imgs: []*image.Summary{
				{
					ID:         "sha256:alpine123",
					Containers: 2,
					Size:       5000000,
					RepoTags:   []string{"alpine:latest"},
				},
			},
			expected: []*protos.DiskUsageImage{
				{
					Id:         "sha256:alpine123",
					Containers: 2,
					Size:       5000000,
					RepoTags:   []string{"alpine:latest"},
				},
			},
		},
		{
			name: "multiple images",
			imgs: []*image.Summary{
				{
					ID:         "sha256:alpine123",
					Containers: 2,
					Size:       5000000,
					RepoTags:   []string{"alpine:latest", "alpine:3.18"},
				},
				{
					ID:         "sha256:nginx456",
					Containers: 1,
					Size:       140000000,
					RepoTags:   []string{"nginx:latest"},
				},
			},
			expected: []*protos.DiskUsageImage{
				{
					Id:         "sha256:alpine123",
					Containers: 2,
					Size:       5000000,
					RepoTags:   []string{"alpine:latest", "alpine:3.18"},
				},
				{
					Id:         "sha256:nginx456",
					Containers: 1,
					Size:       140000000,
					RepoTags:   []string{"nginx:latest"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapDiskUsageImage(tt.imgs)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapDiskUsageContainer(t *testing.T) {
	tests := []struct {
		name     string
		conts    []*container.Summary
		expected []*protos.DiskUsageContainer
	}{
		{
			name:     "empty containers",
			conts:    []*container.Summary{},
			expected: []*protos.DiskUsageContainer{},
		},
		{
			name: "single container",
			conts: []*container.Summary{
				{
					ID:     "abc123",
					Image:  "alpine:latest",
					SizeRw: 1024,
					State:  "running",
					Status: "Up 5 minutes",
				},
			},
			expected: []*protos.DiskUsageContainer{
				{
					Id:     "abc123",
					Image:  "alpine:latest",
					SizeRw: 1024,
					State:  "running",
					Status: "Up 5 minutes",
				},
			},
		},
		{
			name: "multiple containers",
			conts: []*container.Summary{
				{
					ID:     "abc123",
					Image:  "alpine:latest",
					SizeRw: 1024,
					State:  "running",
					Status: "Up 5 minutes",
				},
				{
					ID:     "def456",
					Image:  "nginx:latest",
					SizeRw: 2048,
					State:  "stopped",
					Status: "Exited (0) 1 hour ago",
				},
			},
			expected: []*protos.DiskUsageContainer{
				{
					Id:     "abc123",
					Image:  "alpine:latest",
					SizeRw: 1024,
					State:  "running",
					Status: "Up 5 minutes",
				},
				{
					Id:     "def456",
					Image:  "nginx:latest",
					SizeRw: 2048,
					State:  "stopped",
					Status: "Exited (0) 1 hour ago",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapDiskUsageContainer(tt.conts)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapDiskUsageVolume(t *testing.T) {
	tests := []struct {
		name     string
		vols     []*volume.Volume
		expected []*protos.DiskUsageVolume
	}{
		{
			name:     "empty volumes",
			vols:     []*volume.Volume{},
			expected: []*protos.DiskUsageVolume{},
		},
		{
			name: "single volume",
			vols: []*volume.Volume{
				{
					Name:       "my-volume",
					Mountpoint: "/var/lib/docker/volumes/my-volume/_data",
					UsageData: &volume.UsageData{
						Size: 1048576,
					},
				},
			},
			expected: []*protos.DiskUsageVolume{
				{
					Name:       "my-volume",
					Size:       1048576,
					Mountpoint: "/var/lib/docker/volumes/my-volume/_data",
				},
			},
		},
		{
			name: "multiple volumes",
			vols: []*volume.Volume{
				{
					Name:       "data-volume",
					Mountpoint: "/var/lib/docker/volumes/data-volume/_data",
					UsageData: &volume.UsageData{
						Size: 2097152,
					},
				},
				{
					Name:       "logs-volume",
					Mountpoint: "/var/lib/docker/volumes/logs-volume/_data",
					UsageData: &volume.UsageData{
						Size: 524288,
					},
				},
			},
			expected: []*protos.DiskUsageVolume{
				{
					Name:       "data-volume",
					Size:       2097152,
					Mountpoint: "/var/lib/docker/volumes/data-volume/_data",
				},
				{
					Name:       "logs-volume",
					Size:       524288,
					Mountpoint: "/var/lib/docker/volumes/logs-volume/_data",
				},
			},
		},
		{
			name: "volume without usage data",
			vols: []*volume.Volume{
				{
					Name:       "empty-volume",
					Mountpoint: "/var/lib/docker/volumes/empty-volume/_data",
					UsageData:  nil,
				},
			},
			expected: []*protos.DiskUsageVolume{
				{
					Name:       "empty-volume",
					Size:       0,
					Mountpoint: "/var/lib/docker/volumes/empty-volume/_data",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapDiskUsageVolume(tt.vols)
			assert.Equal(t, tt.expected, result)
		})
	}
}
