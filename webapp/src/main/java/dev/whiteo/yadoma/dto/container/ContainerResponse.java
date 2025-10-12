package dev.whiteo.yadoma.dto.container;

import java.time.LocalDateTime;

/**
 * DTO for container response.
 * Contains container details such as id, name, creation time, status, and state.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
public record ContainerResponse(String id, String name, LocalDateTime createdAt, String status, String state) {}