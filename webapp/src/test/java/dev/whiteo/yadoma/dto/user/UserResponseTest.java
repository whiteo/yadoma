package dev.whiteo.yadoma.dto.user;

import dev.whiteo.yadoma.domain.Role;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class UserResponseTest {

    @Test
    void userResponse_ShouldCreateWithEmailAndRole() {
        String email = "test@example.com";
        Role role = Role.USER;

        UserResponse response = new UserResponse("", email, role);

        assertEquals(email, response.email());
        assertEquals(role, response.role());
    }

    @Test
    void userResponse_ShouldAllowNullValues() {
        UserResponse response = new UserResponse(null,null, null);

        assertNull(response.email());
        assertNull(response.role());
    }

    @Test
    void userResponse_ShouldSupportEquality() {
        UserResponse response1 = new UserResponse("","test@example.com", Role.USER);
        UserResponse response2 = new UserResponse("","test@example.com", Role.USER);
        UserResponse response3 = new UserResponse("","test@example.com", Role.ADMIN);

        assertEquals(response1, response2);
        assertNotEquals(response1, response3);
        assertEquals(response1.hashCode(), response2.hashCode());
    }

    @Test
    void userResponse_ShouldHaveProperToString() {
        UserResponse response = new UserResponse("","test@example.com", Role.USER);

        String toString = response.toString();

        assertNotNull(toString);
        assertTrue(toString.contains("test@example.com"));
        assertTrue(toString.contains("USER"));
    }

    @Test
    void userResponse_ShouldWorkWithDifferentRoles() {
        UserResponse userResponse = new UserResponse("","user@example.com", Role.USER);
        UserResponse adminResponse = new UserResponse("","admin@example.com", Role.ADMIN);

        assertEquals(Role.USER, userResponse.role());
        assertEquals(Role.ADMIN, adminResponse.role());
        assertNotEquals(userResponse, adminResponse);
    }
}
