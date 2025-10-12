package dev.whiteo.yadoma.controller;

import dev.whiteo.yadoma.dto.user.UserCreateRequest;
import dev.whiteo.yadoma.dto.user.UserResponse;
import dev.whiteo.yadoma.dto.user.UserUpdatePasswordRequest;
import dev.whiteo.yadoma.exception.ResponseError;
import dev.whiteo.yadoma.security.AuthInterceptor;
import dev.whiteo.yadoma.security.AuthRequired;
import dev.whiteo.yadoma.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

/**
 * Controller for user-related endpoints.
 * Provides CRUD operations for user management.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/user")
@AuthRequired(AuthInterceptor.class)
@Tag(name = "user", description = "User endpoints")
public class UserRestController {

    private final AuthInterceptor authInterceptor;
    private final UserService service;

    /**
     * Retrieves the details of all users.
     *
     * @return ResponseEntity containing the list of user details or error response
     */
    @GetMapping("/all")
    @SecurityRequirement(name = "Bearer Authentication")
    @Operation(summary = "Finds all users")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "User details",
                    content = @Content(schema = @Schema(implementation = UserResponse.class))),
            @ApiResponse(responseCode = "401", description = "Unauthorized",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<List<UserResponse>> findAll() {
        return ResponseEntity.ok(service.findAll(authInterceptor.getUserId()));
    }

    /**
     * Retrieves the details of the currently authenticated user.
     *
     * @return ResponseEntity containing the user details or error response
     */
    @GetMapping("/me")
    @SecurityRequirement(name = "Bearer Authentication")
    @Operation(summary = "Finds the current user")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "User details",
                    content = @Content(schema = @Schema(implementation = UserResponse.class))),
            @ApiResponse(responseCode = "401", description = "Unauthorized",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<UserResponse> me() {
        return ResponseEntity.ok(service.me(authInterceptor.getUserId()));
    }

    /**
     * Creates a new user.
     *
     * @param request User creation request
     * @return ResponseEntity containing the created user or error response
     */
    @PostMapping("/create")
    @Operation(summary = "Creates a new user")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "201", description = "User created successfully"),
            @ApiResponse(responseCode = "409", description = "Conflict - user already exists",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<Void> create(@Valid @RequestBody UserCreateRequest request) {
        service.create(request);
        return ResponseEntity.status(HttpStatus.CREATED).build();
    }

    /**
     * Updates the password of the currently authenticated user.
     *
     * @param request Password update request
     * @return ResponseEntity indicating success or error
     */
    @PostMapping("/update-password")
    @SecurityRequirement(name = "Bearer Authentication")
    @Operation(summary = "Updates the password of the current user")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Password updated successfully"),
            @ApiResponse(responseCode = "401", description = "Unauthorized",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<Void> updatePassword(@Valid @RequestBody UserUpdatePasswordRequest request) {
        service.updatePassword(authInterceptor.getUserId(), request);
        return ResponseEntity.ok().build();
    }

    /**
     * Deletes the currently authenticated user.
     *
     * @return ResponseEntity indicating success or error
     */
    @DeleteMapping("/delete/{id}")
    @SecurityRequirement(name = "Bearer Authentication")
    @Operation(summary = "Deletes the current user")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "204", description = "User deleted successfully"),
            @ApiResponse(responseCode = "401", description = "Unauthorized",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "500", description = "Internal server error",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
            @ApiResponse(responseCode = "503", description = "Service unavailable",
                    content = @Content(schema = @Schema(implementation = ResponseError.class))),
    })
    public ResponseEntity<Void> delete(@PathVariable("id") String deleteId) {
        service.remove(deleteId,authInterceptor.getUserId());
        return ResponseEntity.noContent().build();
    }
}