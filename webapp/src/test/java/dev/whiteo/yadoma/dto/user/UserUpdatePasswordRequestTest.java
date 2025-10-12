package dev.whiteo.yadoma.dto.user;

import jakarta.validation.ConstraintViolation;
import jakarta.validation.Validation;
import jakarta.validation.Validator;
import jakarta.validation.ValidatorFactory;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Set;

import static org.junit.jupiter.api.Assertions.*;

class UserUpdatePasswordRequestTest {

    private Validator validator;

    @BeforeEach
    void setUp() {
        ValidatorFactory factory = Validation.buildDefaultValidatorFactory();
        validator = factory.getValidator();
    }

    @Test
    void validUserUpdatePasswordRequest_ShouldPassValidation() {
        // Given
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest("newPassword123", "oldPassword123");

        // When
        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        // Then
        assertTrue(violations.isEmpty());
        assertEquals("newPassword123", request.password());
        assertEquals("oldPassword123", request.oldPassword());
    }

    @Test
    void blankPassword_ShouldFailValidation() {
        // Given
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest("", "oldPassword123");

        // When
        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Password cannot be blank")));
    }

    @Test
    void nullPassword_ShouldFailValidation() {
        // Given
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest(null, "oldPassword123");

        // When
        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Password cannot be blank")));
    }

    @Test
    void shortPassword_ShouldFailValidation() {
        // Given
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest("1234", "oldPassword123");

        // When
        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Password must be between 5 and 30 characters")));
    }

    @Test
    void longPassword_ShouldFailValidation() {
        // Given
        String longPassword = "a".repeat(31);
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest(longPassword, "oldPassword123");

        // When
        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Password must be between 5 and 30 characters")));
    }

    @Test
    void blankOldPassword_ShouldFailValidation() {
        // Given
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest("newPassword123", "");

        // When
        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Old password cannot be blank")));
    }

    @Test
    void nullOldPassword_ShouldFailValidation() {
        // Given
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest("newPassword123", null);

        // When
        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Old password cannot be blank")));
    }

    @Test
    void shortOldPassword_ShouldFailValidation() {
        // Given
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest("newPassword123", "1234");

        // When
        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Old password must be between 5 and 30 characters")));
    }

    @Test
    void longOldPassword_ShouldFailValidation() {
        // Given
        String longOldPassword = "a".repeat(31);
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest("newPassword123", longOldPassword);

        // When
        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Old password must be between 5 and 30 characters")));
    }
}
