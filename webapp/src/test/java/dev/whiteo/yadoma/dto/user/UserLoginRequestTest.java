package dev.whiteo.yadoma.dto.user;

import jakarta.validation.ConstraintViolation;
import jakarta.validation.Validation;
import jakarta.validation.Validator;
import jakarta.validation.ValidatorFactory;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Set;

import static org.junit.jupiter.api.Assertions.*;

class UserLoginRequestTest {

    private Validator validator;

    @BeforeEach
    void setUp() {
        ValidatorFactory factory = Validation.buildDefaultValidatorFactory();
        validator = factory.getValidator();
    }

    @Test
    void validUserLoginRequest_ShouldPassValidation() {

        UserLoginRequest request = new UserLoginRequest("test@example.com", "password123");

        Set<ConstraintViolation<UserLoginRequest>> violations = validator.validate(request);

        assertTrue(violations.isEmpty());
        assertEquals("test@example.com", request.email());
        assertEquals("password123", request.password());
    }

    @Test
    void invalidEmailFormat_ShouldFailValidation() {
        UserLoginRequest request = new UserLoginRequest("invalid-email", "password123");

        Set<ConstraintViolation<UserLoginRequest>> violations = validator.validate(request);

        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Invalid email format")));
    }

    @Test
    void blankEmail_ShouldFailValidation() {
        UserLoginRequest request = new UserLoginRequest("", "password123");

        Set<ConstraintViolation<UserLoginRequest>> violations = validator.validate(request);

        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Email cannot be blank")));
    }

    @Test
    void nullEmail_ShouldFailValidation() {
        UserLoginRequest request = new UserLoginRequest(null, "password123");

        Set<ConstraintViolation<UserLoginRequest>> violations = validator.validate(request);

        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Email cannot be blank")));
    }

    @Test
    void blankPassword_ShouldFailValidation() {
        UserLoginRequest request = new UserLoginRequest("test@example.com", "");

        Set<ConstraintViolation<UserLoginRequest>> violations = validator.validate(request);

        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Password cannot be blank")));
    }

    @Test
    void nullPassword_ShouldFailValidation() {
        UserLoginRequest request = new UserLoginRequest("test@example.com", null);

        Set<ConstraintViolation<UserLoginRequest>> violations = validator.validate(request);

        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Password cannot be blank")));
    }

    @Test
    void shortPassword_ShouldFailValidation() {
        UserLoginRequest request = new UserLoginRequest("test@example.com", "1234");

        Set<ConstraintViolation<UserLoginRequest>> violations = validator.validate(request);

        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Password must be between 5 and 30 characters")));
    }

    @Test
    void longPassword_ShouldFailValidation() {
        String longPassword = "a".repeat(31);
        UserLoginRequest request = new UserLoginRequest("test@example.com", longPassword);

        Set<ConstraintViolation<UserLoginRequest>> violations = validator.validate(request);

        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Password must be between 5 and 30 characters")));
    }
}
