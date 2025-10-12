package dev.whiteo.yadoma.dto.token;

/**
 * DTO for user login response.
 * Contains JWT token and user information.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
public record TokenResponse(String token) {}