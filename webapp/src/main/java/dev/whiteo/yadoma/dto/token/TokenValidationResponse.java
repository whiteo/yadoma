package dev.whiteo.yadoma.dto.token;

/**
 * Represents the response for token validation.
 * Contains a boolean indicating whether the token is valid.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
public record TokenValidationResponse(Boolean valid) {}