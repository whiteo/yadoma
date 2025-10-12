package dev.whiteo.yadoma.controller;

import container.v1.Container;
import dev.whiteo.yadoma.dto.container.ContainerCreateRequest;
import dev.whiteo.yadoma.dto.container.ContainerResponse;
import dev.whiteo.yadoma.exception.ResponseError;
import dev.whiteo.yadoma.security.AuthInterceptor;
import dev.whiteo.yadoma.security.AuthRequired;
import dev.whiteo.yadoma.service.ContainerService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.servlet.mvc.method.annotation.StreamingResponseBody;

import java.util.Iterator;
import java.util.List;

/**
 * Controller for container endpoints.
 * Provides CRUD operations for container management.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/container")
@AuthRequired(AuthInterceptor.class)
@Tag(name = "container", description = "Container endpoints")
public class ContainerRestController {

    private final AuthInterceptor authInterceptor;
    private final ContainerService service;

    /**
     * Retrieves a list of containers for the authenticated user.
     *
     * @param userId The ID of the authenticated user.
     * @return ResponseEntity containing a list of ContainerResponse objects.
     */
    @GetMapping("/{userId}/all")
    @SecurityRequirement(name = "Bearer Authentication")
    @Operation(summary = "Finds all containers for the user")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Containers details",
                    content = @Content(schema = @Schema(implementation = ContainerResponse.class))),
            @ApiResponse(responseCode = "401", description = "Unauthorized",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<List<ContainerResponse>> findAll(@PathVariable("userId") String userId) {
        return ResponseEntity.ok(service.findAll(userId, authInterceptor.getUserId()));
    }

    /**
     * Retrieves a container by its ID for the authenticated user.
     *
     * @param containerId The ID of the container to retrieve.
     * @return ResponseEntity containing a ContainerResponse object.
     */
    @GetMapping("/get/{containerId}")
    @SecurityRequirement(name = "Bearer Authentication")
    @Operation(summary = "Finds a container for the user")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Container details",
                    content = @Content(schema = @Schema(implementation = ContainerResponse.class))),
            @ApiResponse(responseCode = "401", description = "Unauthorized",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<ContainerResponse> getById(@PathVariable("containerId") String containerId) {
        return ResponseEntity.ok(service.getById(containerId, authInterceptor.getUserId()));
    }

    /**
     * Creates a new container for the authenticated user.
     *
     * @param request The ContainerCreateRequest containing container details.
     * @return ResponseEntity indicating the success or failure of the operation.
     */
    @PostMapping("/create")
    @SecurityRequirement(name = "Bearer Authentication")
    @Operation(summary = "Creates a new container for the user")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Container created successfully"),
            @ApiResponse(responseCode = "401", description = "Unauthorized",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<Void> create(@Valid @RequestBody ContainerCreateRequest request) {
        service.create(request, authInterceptor.getUserId());
        return ResponseEntity.ok().build();
    }

    /**
     * Deletes a container by its ID for the authenticated user.
     *
     * @param containerId The ID of the container to delete.
     * @return ResponseEntity indicating the success or failure of the operation.
     */
    @DeleteMapping("/delete/{containerId}")
    @SecurityRequirement(name = "Bearer Authentication")
    @Operation(summary = "Deletes a container for the user")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "204", description = "Container deleted successfully"),
            @ApiResponse(responseCode = "401", description = "Unauthorized",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<Void> delete(@PathVariable("containerId") String containerId) {
        service.delete(containerId, authInterceptor.getUserId());
        return ResponseEntity.noContent().build();
    }

    /**
     * Starts a container by its ID for the authenticated user.
     *
     * @param containerId The ID of the container to start.
     * @return ResponseEntity indicating the success or failure of the operation.
     */
    @PostMapping("/start/{containerId}")
    @SecurityRequirement(name = "Bearer Authentication")
    @Operation(summary = "Starts a container for the user")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Container started successfully"),
            @ApiResponse(responseCode = "401", description = "Unauthorized",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "409", description = "Conflict - container already running",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<Void> start(@PathVariable("containerId") String containerId) {
        service.start(containerId, authInterceptor.getUserId());
        return ResponseEntity.ok().build();
    }

    /**
     * Stops a container by its ID for the authenticated user.
     *
     * @param containerId The ID of the container to stop.
     * @return ResponseEntity indicating the success or failure of the operation.
     */
    @PostMapping("/stop/{containerId}")
    @SecurityRequirement(name = "Bearer Authentication")
    @Operation(summary = "Stops a container for the user")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Container stopped successfully"),
            @ApiResponse(responseCode = "401", description = "Unauthorized",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "409", description = "Conflict - container is not running",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<Void> stop(@PathVariable("containerId") String containerId) {
        service.stop(containerId, authInterceptor.getUserId());
        return ResponseEntity.ok().build();
    }

    /**
     * Restarts a container by its ID for the authenticated user.
     *
     * @param containerId The ID of the container to restart.
     * @return ResponseEntity indicating the success or failure of the operation.
     */
    @PostMapping("/restart/{containerId}")
    @SecurityRequirement(name = "Bearer Authentication")
    @Operation(summary = "Restarts a container for the user")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Container restarted successfully"),
            @ApiResponse(responseCode = "401", description = "Unauthorized",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "409", description = "Conflict - container is not running",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<Void> restart(@PathVariable("containerId") String containerId) {
        service.restart(containerId, authInterceptor.getUserId());
        return ResponseEntity.ok().build();
    }

    /**
     * Retrieves logs for a container by its ID.
     *
     * @param containerId The ID of the container to get logs from.
     * @param follow      Whether to follow the log stream (default: false).
     * @return ResponseEntity with streaming log content.
     */
    @GetMapping(value = "/logs/{containerId}", produces = MediaType.TEXT_PLAIN_VALUE)
    @SecurityRequirement(name = "Bearer Authentication")
    @Operation(summary = "Get container logs")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Container logs retrieved successfully"),
            @ApiResponse(responseCode = "401", description = "Unauthorized",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "404", description = "Container not found",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<StreamingResponseBody> getLogs(
            @PathVariable("containerId") String containerId,
            @RequestParam(defaultValue = "false") boolean follow) {

        final String userId = authInterceptor.getUserId();

        StreamingResponseBody stream = outputStream -> {
            try {
                Iterator<Container.GetContainerLogsResponse> logs =
                        service.getLogs(containerId, userId, follow);

                while (logs.hasNext()) {
                    Container.GetContainerLogsResponse chunk = logs.next();
                    outputStream.write(chunk.getChunk().toByteArray());
                    outputStream.flush();
                }
            } catch (Exception e) {
                outputStream.write(("Error retrieving logs: " + e.getMessage()).getBytes());
                outputStream.flush();
            }
        };

        return ResponseEntity.ok()
                .contentType(MediaType.TEXT_PLAIN)
                .body(stream);
    }
}