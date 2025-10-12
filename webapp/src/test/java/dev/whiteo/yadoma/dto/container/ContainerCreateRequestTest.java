package dev.whiteo.yadoma.dto.container;

import jakarta.validation.ConstraintViolation;
import jakarta.validation.Validation;
import jakarta.validation.Validator;
import jakarta.validation.ValidatorFactory;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.Set;

import static org.junit.jupiter.api.Assertions.*;

class ContainerCreateRequestTest {

    private Validator validator;

    @BeforeEach
    void setUp() {
        ValidatorFactory factory = Validation.buildDefaultValidatorFactory();
        validator = factory.getValidator();
    }

    @Test
    void validContainerCreateRequest_ShouldPassValidation() {
        // Given
        List<String> envVars = Arrays.asList("ENV_VAR1=value1", "ENV_VAR2=value2");
        ContainerCreateRequest request = new ContainerCreateRequest("test-container", "nginx:latest", envVars);

        // When
        Set<ConstraintViolation<ContainerCreateRequest>> violations = validator.validate(request);

        // Then
        assertTrue(violations.isEmpty());
        assertEquals("test-container", request.name());
        assertEquals("nginx:latest", request.image());
        assertEquals(envVars, request.envVars());
    }

    @Test
    void blankName_ShouldFailValidation() {
        // Given
        ContainerCreateRequest request = new ContainerCreateRequest("", "nginx:latest", Collections.emptyList());

        // When
        Set<ConstraintViolation<ContainerCreateRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Name cannot be blank")));
    }

    @Test
    void nullName_ShouldFailValidation() {
        // Given
        ContainerCreateRequest request = new ContainerCreateRequest(null, "nginx:latest", Collections.emptyList());

        // When
        Set<ConstraintViolation<ContainerCreateRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Name cannot be blank")));
    }

    @Test
    void blankImage_ShouldFailValidation() {
        // Given
        ContainerCreateRequest request = new ContainerCreateRequest("test-container", "", Collections.emptyList());

        // When
        Set<ConstraintViolation<ContainerCreateRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Image cannot be blank")));
    }

    @Test
    void nullImage_ShouldFailValidation() {
        // Given
        ContainerCreateRequest request = new ContainerCreateRequest("test-container", null, Collections.emptyList());

        // When
        Set<ConstraintViolation<ContainerCreateRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Image cannot be blank")));
    }

    @Test
    void nullEnvVars_ShouldFailValidation() {
        // Given
        ContainerCreateRequest request = new ContainerCreateRequest("test-container", "nginx:latest", null);

        // When
        Set<ConstraintViolation<ContainerCreateRequest>> violations = validator.validate(request);

        // Then
        assertFalse(violations.isEmpty());
        assertTrue(violations.stream().anyMatch(v -> v.getMessage().contains("Environment variables cannot be null")));
    }

    @Test
    void emptyEnvVars_ShouldPassValidation() {
        // Given
        ContainerCreateRequest request = new ContainerCreateRequest("test-container", "nginx:latest", Collections.emptyList());

        // When
        Set<ConstraintViolation<ContainerCreateRequest>> violations = validator.validate(request);

        // Then
        assertTrue(violations.isEmpty());
    }
}
