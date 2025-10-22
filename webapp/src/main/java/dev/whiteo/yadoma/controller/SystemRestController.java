package dev.whiteo.yadoma.controller;

import dev.whiteo.yadoma.dto.system.DiskUsageResponse;
import dev.whiteo.yadoma.dto.system.SystemInfoResponse;
import dev.whiteo.yadoma.exception.ResponseError;
import dev.whiteo.yadoma.security.AuthInterceptor;
import dev.whiteo.yadoma.security.AuthRequired;
import dev.whiteo.yadoma.service.SystemService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * REST controller for system information endpoints.
 * Provides endpoints to retrieve Docker system info and disk usage.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/system")
@AuthRequired(AuthInterceptor.class)
@Tag(name = "system", description = "System information endpoints")
public class SystemRestController {

    private final SystemService service;

    /**
     * Retrieves system information from Docker daemon.
     *
     * @return ResponseEntity containing system information
     */
    @GetMapping("/info")
    @SecurityRequirement(name = "Bearer Authentication")
    @Operation(summary = "Get Docker system information")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "System information retrieved successfully",
                    content = @Content(schema = @Schema(implementation = SystemInfoResponse.class))),
            @ApiResponse(responseCode = "401", description = "Unauthorized",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<SystemInfoResponse> getSystemInfo() {
        return ResponseEntity.ok(service.getSystemInfo());
    }

    /**
     * Retrieves disk usage information for images, containers, and volumes.
     *
     * @return ResponseEntity containing disk usage information
     */
    @GetMapping("/disk-usage")
    @SecurityRequirement(name = "Bearer Authentication")
    @Operation(summary = "Get disk usage information")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Disk usage retrieved successfully",
                    content = @Content(schema = @Schema(implementation = DiskUsageResponse.class))),
            @ApiResponse(responseCode = "401", description = "Unauthorized",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<DiskUsageResponse> getDiskUsage() {
        return ResponseEntity.ok(service.getDiskUsage());
    }
}
