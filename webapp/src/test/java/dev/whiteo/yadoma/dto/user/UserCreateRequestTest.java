package dev.whiteo.yadoma.dto.user;

import jakarta.validation.ConstraintViolation;
import jakarta.validation.Validation;
import jakarta.validation.Validator;
import jakarta.validation.ValidatorFactory;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Set;

import static org.junit.jupiter.api.Assertions.*;

class UserCreateRequestTest {

    private Validator validator;

    @BeforeEach
    void setUp() {
        ValidatorFactory factory = Validation.buildDefaultValidatorFactory();
        validator = factory.getValidator();
    }

    @Test
    void validUserCreateRequest_ShouldPassValidation() {
        // Given
        UserCreateRequest request = new UserCreateRequest("test@example.com", "password123");

        // When
        Set<ConstraintViolation<UserCreateRequest>> violations = validator.validate(request);

        // Then
        assertTrue(violations.isEmpty());
        assertEquals("test@example.com", request.email());
        assertEquals("password123", request.password());
    }

    @Test
    void invalidEmailFormat_ShouldFailValidation() {
        // Given
        UserCreateRequest request = new UserCreateRequest("invalid-email", "password123");

        // When
        Set<ConstraintViolation<UserCreateRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Invalid email format")));
    }

    @Test
    void blankEmail_ShouldFailValidation() {
        // Given
        UserCreateRequest request = new UserCreateRequest("", "password123");

        // When
        Set<ConstraintViolation<UserCreateRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Email cannot be blank")));
    }

    @Test
    void nullEmail_ShouldFailValidation() {
        // Given
        UserCreateRequest request = new UserCreateRequest(null, "password123");

        // When
        Set<ConstraintViolation<UserCreateRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Email cannot be blank")));
    }

    @Test
    void blankPassword_ShouldFailValidation() {
        // Given
        UserCreateRequest request = new UserCreateRequest("test@example.com", "");

        // When
        Set<ConstraintViolation<UserCreateRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Password cannot be blank")));
    }

    @Test
    void nullPassword_ShouldFailValidation() {
        // Given
        UserCreateRequest request = new UserCreateRequest("test@example.com", null);

        // When
        Set<ConstraintViolation<UserCreateRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Password cannot be blank")));
    }

    @Test
    void shortPassword_ShouldFailValidation() {
        // Given
        UserCreateRequest request = new UserCreateRequest("test@example.com", "1234");

        // When
        Set<ConstraintViolation<UserCreateRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Password must be between 5 and 30 characters")));
    }

    @Test
    void longPassword_ShouldFailValidation() {
        // Given
        String longPassword = "a".repeat(31);
        UserCreateRequest request = new UserCreateRequest("test@example.com", longPassword);

        // When
        Set<ConstraintViolation<UserCreateRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Password must be between 5 and 30 characters")));
    }

    @Test
    void equalObjects_ShouldBeEqual() {
        // Given
        UserCreateRequest request1 = new UserCreateRequest("test@example.com", "password123");
        UserCreateRequest request2 = new UserCreateRequest("test@example.com", "password123");

        // Then
        assertEquals(request1, request2);
        assertEquals(request1.hashCode(), request2.hashCode());
    }

    @Test
    void differentObjects_ShouldNotBeEqual() {
        // Given
        UserCreateRequest request1 = new UserCreateRequest("test1@example.com", "password123");
        UserCreateRequest request2 = new UserCreateRequest("test2@example.com", "password123");

        // Then
        assertNotEquals(request1, request2);
    }
}
