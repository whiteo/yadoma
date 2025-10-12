package dev.whiteo.yadoma.dto.user;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

/**
 * DTO for updating user's password.
 * Contains new and old password fields with validation constraints.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
public record UserUpdatePasswordRequest(
        @NotBlank(message = "Password cannot be blank")
        @Size(min = 5, max = 30, message = "Password must be between 5 and 30 characters")
        String password,


        @NotBlank(message = "Old password cannot be blank")
        @Size(min = 5, max = 30, message = "Old password must be between 5 and 30 characters")
        String oldPassword
) {}