package dev.whiteo.yadoma.dto.user;

import dev.whiteo.yadoma.domain.Role;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class UserResponseTest {

    @Test
    void userResponse_ShouldCreateWithEmailAndRole() {
        // Given
        String email = "test@example.com";
        Role role = Role.USER;

        // When
        UserResponse response = new UserResponse(email, role);

        // Then
        assertEquals(email, response.email());
        assertEquals(role, response.role());
    }

    @Test
    void userResponse_ShouldAllowNullValues() {
        // When
        UserResponse response = new UserResponse(null, null);

        // Then
        assertNull(response.email());
        assertNull(response.role());
    }

    @Test
    void userResponse_ShouldSupportEquality() {
        // Given
        UserResponse response1 = new UserResponse("test@example.com", Role.USER);
        UserResponse response2 = new UserResponse("test@example.com", Role.USER);
        UserResponse response3 = new UserResponse("test@example.com", Role.ADMIN);

        // Then
        assertEquals(response1, response2);
        assertNotEquals(response1, response3);
        assertEquals(response1.hashCode(), response2.hashCode());
    }

    @Test
    void userResponse_ShouldHaveProperToString() {
        // Given
        UserResponse response = new UserResponse("test@example.com", Role.USER);

        // When
        String toString = response.toString();

        // Then
        assertNotNull(toString);
        assertTrue(toString.contains("test@example.com"));
        assertTrue(toString.contains("USER"));
    }

    @Test
    void userResponse_ShouldWorkWithDifferentRoles() {
        // Given
        UserResponse userResponse = new UserResponse("user@example.com", Role.USER);
        UserResponse adminResponse = new UserResponse("admin@example.com", Role.ADMIN);

        // Then
        assertEquals(Role.USER, userResponse.role());
        assertEquals(Role.ADMIN, adminResponse.role());
        assertNotEquals(userResponse, adminResponse);
    }
}
