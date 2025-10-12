package dev.whiteo.yadoma.dto.container;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;

import java.util.List;

/**
 * Represents a container creation request DTO.
 * Contains name, image, and environment variables.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
public record ContainerCreateRequest(
        @NotBlank(message = "Name cannot be blank")
        String name,

        @NotBlank(message = "Image cannot be blank")
        String image,

        @NotNull(message = "Environment variables cannot be null")
        List<String> envVars
) {}