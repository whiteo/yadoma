package dev.whiteo.yadoma.controller;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * REST controller for providing application version information.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@RestController
@RequiredArgsConstructor
@RequestMapping("/version")
@Tag(name = "version", description = "Version endpoints")
public class VersionRestController {

    @Value("${app.version}")
    private String appVersion;

    /**
     * Returns the current version of the application.
     * @return application version string
     */
    @GetMapping
    @Operation(summary = "Get application version")
    @ApiResponse(responseCode = "200", description = "Application version")
    public String getAppVersion() {
        return appVersion;
    }
}