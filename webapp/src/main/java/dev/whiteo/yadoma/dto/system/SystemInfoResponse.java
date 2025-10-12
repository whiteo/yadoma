package dev.whiteo.yadoma.dto.system;

import java.util.List;

/**
 * DTO for system information response.
 * Contains Docker system details, resource usage, and statistics.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
public record SystemInfoResponse(
        String id,
        String name,
        String serverVersion,
        String kernelVersion,
        String operatingSystem,
        String architecture,
        Integer nCpu,
        Long memTotal,
        Integer containers,
        Integer containersRunning,
        Integer containersPaused,
        Integer containersStopped,
        Integer images,
        String driver,
        List<String> labels
) {}
