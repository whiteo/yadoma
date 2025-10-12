package dev.whiteo.yadoma.dto.system;

import java.util.List;

/**
 * DTO for disk usage image information.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
public record DiskUsageImageResponse(
        String id,
        List<String> repoTags,
        Long size,
        Long containers
) {}
