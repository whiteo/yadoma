package dev.whiteo.yadoma.dto.system;

/**
 * DTO for disk usage container information.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
public record DiskUsageContainerResponse(
        String id,
        String image,
        String state,
        String status,
        Long sizeRw
) {}
