package dev.whiteo.yadoma.controller;

import dev.whiteo.yadoma.dto.token.TokenResponse;
import dev.whiteo.yadoma.dto.token.TokenValidationResponse;
import dev.whiteo.yadoma.dto.user.UserLoginRequest;
import dev.whiteo.yadoma.exception.ResponseError;
import dev.whiteo.yadoma.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * REST controller for user authentication and token validation.
 * Provides endpoints for login and token verification.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/authenticate")
@Tag(name = "authenticate", description = "Authentication endpoints")
public class AuthenticateRestController {

    private final UserService service;

    /**
     * Authenticates a user and returns a JWT token.
     *
     * @param request user login request
     * @return JWT token and user info
     */
    @PostMapping
    @Operation(summary = "Authenticate a user and return a JWT token")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "JWT token and user info",
                    content = @Content(schema = @Schema(implementation = TokenResponse.class))),
            @ApiResponse(responseCode = "400", description = "Invalid username or password",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<TokenResponse> authenticate(@Valid @RequestBody UserLoginRequest request) {
        return ResponseEntity.ok(service.getToken(request));
    }

    /**
     * Validates the JWT token from the request.
     *
     * @return true if token is valid, false otherwise
     */
    @GetMapping
    @Operation(summary = "Validate JWT token")
    @SecurityRequirement(name = "Bearer Authentication")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "True if token is valid, false otherwise",
                    content = @Content(schema = @Schema(implementation = TokenValidationResponse.class))),
            @ApiResponse(responseCode = "401", description = "Unauthorized",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<TokenValidationResponse> validate(HttpServletRequest request) {
        return ResponseEntity.ok(service.validateToken(request));
    }
}