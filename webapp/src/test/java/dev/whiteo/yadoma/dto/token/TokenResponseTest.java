package dev.whiteo.yadoma.dto.token;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class TokenResponseTest {

    @Test
    void tokenResponse_ShouldCreateWithToken() {
        // Given
        String token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...";

        // When
        TokenResponse response = new TokenResponse(token);

        // Then
        assertEquals(token, response.token());
    }

    @Test
    void tokenResponse_ShouldAllowNullToken() {
        // When
        TokenResponse response = new TokenResponse(null);

        // Then
        assertNull(response.token());
    }

    @Test
    void tokenResponse_ShouldSupportEquality() {
        // Given
        TokenResponse response1 = new TokenResponse("token123");
        TokenResponse response2 = new TokenResponse("token123");
        TokenResponse response3 = new TokenResponse("token456");

        // Then
        assertEquals(response1, response2);
        assertNotEquals(response1, response3);
        assertEquals(response1.hashCode(), response2.hashCode());
    }

    @Test
    void tokenResponse_ShouldHaveProperToString() {
        // Given
        TokenResponse response = new TokenResponse("token123");

        // When
        String toString = response.toString();

        // Then
        assertNotNull(toString);
        assertTrue(toString.contains("token123"));
    }
}
