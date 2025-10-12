package dev.whiteo.yadoma.dto.user;

import dev.whiteo.yadoma.domain.Role;

/**
 * DTO for user response.
 * Contains user's email and role information.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
public record UserResponse(String email, Role role) {}