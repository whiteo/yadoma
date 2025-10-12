package dev.whiteo.yadoma.dto.token;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class TokenResponseTest {

    @Test
    void tokenResponse_ShouldCreateWithToken() {
        String token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...";

        TokenResponse response = new TokenResponse(token);

        assertEquals(token, response.token());
    }

    @Test
    void tokenResponse_ShouldAllowNullToken() {
        TokenResponse response = new TokenResponse(null);

        assertNull(response.token());
    }

    @Test
    void tokenResponse_ShouldSupportEquality() {
        TokenResponse response1 = new TokenResponse("token123");
        TokenResponse response2 = new TokenResponse("token123");
        TokenResponse response3 = new TokenResponse("token456");

        assertEquals(response1, response2);
        assertNotEquals(response1, response3);
        assertEquals(response1.hashCode(), response2.hashCode());
    }

    @Test
    void tokenResponse_ShouldHaveProperToString() {
        TokenResponse response = new TokenResponse("token123");

        String toString = response.toString();

        assertNotNull(toString);
        assertTrue(toString.contains("token123"));
    }
}
