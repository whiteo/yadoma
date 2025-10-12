package dev.whiteo.yadoma.dto.system;

import java.util.List;

/**
 * DTO for disk usage response.
 * Contains comprehensive disk usage information for images, containers, and volumes.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
public record DiskUsageResponse(
        Long layersSize,
        List<DiskUsageImageResponse> images,
        List<DiskUsageContainerResponse> containers,
        List<DiskUsageVolumeResponse> volumes
) {}
