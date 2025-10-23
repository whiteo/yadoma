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
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest("newPassword123", "oldPassword123");

        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        assertTrue(violations.isEmpty());
        assertEquals("newPassword123", request.password());
        assertEquals("oldPassword123", request.oldPassword());
    }

    @Test
    void blankPassword_ShouldFailValidation() {
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest("", "oldPassword123");

        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Password cannot be blank")));
    }

    @Test
    void nullPassword_ShouldFailValidation() {
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest(null, "oldPassword123");

        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Password cannot be blank")));
    }

    @Test
    void shortPassword_ShouldFailValidation() {
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest("1234", "oldPassword123");

        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Password must be between 5 and 30 characters")));
    }

    @Test
    void longPassword_ShouldFailValidation() {
        String longPassword = "a".repeat(31);
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest(longPassword, "oldPassword123");

        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Password must be between 5 and 30 characters")));
    }

    @Test
    void blankOldPassword_ShouldFailValidation() {
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest("newPassword123", "");

        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Old password cannot be blank")));
    }

    @Test
    void nullOldPassword_ShouldFailValidation() {
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest("newPassword123", null);

        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Old password cannot be blank")));
    }

    @Test
    void shortOldPassword_ShouldFailValidation() {
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest("newPassword123", "1234");

        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Old password must be between 5 and 30 characters")));
    }

    @Test
    void longOldPassword_ShouldFailValidation() {
        String longOldPassword = "a".repeat(31);
        UserUpdatePasswordRequest request = new UserUpdatePasswordRequest("newPassword123", longOldPassword);

        Set<ConstraintViolation<UserUpdatePasswordRequest>> violations = validator.validate(request);

        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Old password must be between 5 and 30 characters")));
    }
}
