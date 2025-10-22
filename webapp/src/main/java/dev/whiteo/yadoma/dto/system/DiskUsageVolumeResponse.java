package dev.whiteo.yadoma.dto.system;

/**
 * DTO for disk usage volume information.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
public record DiskUsageVolumeResponse(
        String name,
        String mountpoint,
        Long size
) {}
