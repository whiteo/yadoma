// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package system

import (
	"github.com/whiteo/yadoma/internal/protos"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/volume"

	"github.com/docker/docker/api/types/image"
)

func mapDiskUsageImage(imgs []*image.Summary) []*protos.DiskUsageImage {
	r := make([]*protos.DiskUsageImage, 0, len(imgs))

	for _, img := range imgs {
		r = append(r, &protos.DiskUsageImage{
			Containers: img.Containers,
			Size:       img.Size,
			Id:         img.ID,
			RepoTags:   img.RepoTags,
		})
	}

	return r
}

func mapDiskUsageContainer(conts []*container.Summary) []*protos.DiskUsageContainer {
	r := make([]*protos.DiskUsageContainer, 0, len(conts))

	for _, cont := range conts {
		r = append(r, &protos.DiskUsageContainer{
			Id:     cont.ID,
			Image:  cont.Image,
			SizeRw: cont.SizeRw,
			State:  cont.State,
			Status: cont.Status,
		})
	}

	return r
}

func mapDiskUsageVolume(vols []*volume.Volume) []*protos.DiskUsageVolume {
	r := make([]*protos.DiskUsageVolume, 0, len(vols))

	for _, vol := range vols {
		var size int64
		if vol.UsageData != nil {
			size = vol.UsageData.Size
		}

		r = append(r, &protos.DiskUsageVolume{
			Name:       vol.Name,
			Size:       size,
			Mountpoint: vol.Mountpoint,
		})
	}

	return r
}
