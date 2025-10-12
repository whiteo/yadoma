package dev.whiteo.yadoma.dto.token;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class TokenValidationResponseTest {

    @Test
    void tokenValidationResponse_ShouldCreateWithValidTrue() {
        Boolean valid = true;

        TokenValidationResponse response = new TokenValidationResponse(valid);

        assertEquals(valid, response.valid());
        assertTrue(response.valid());
    }

    @Test
    void tokenValidationResponse_ShouldCreateWithValidFalse() {
        Boolean valid = false;

        TokenValidationResponse response = new TokenValidationResponse(valid);

        assertEquals(valid, response.valid());
        assertFalse(response.valid());
    }

    @Test
    void tokenValidationResponse_ShouldAllowNullValue() {
        TokenValidationResponse response = new TokenValidationResponse(null);

        assertNull(response.valid());
    }

    @Test
    void tokenValidationResponse_ShouldSupportEquality() {
        TokenValidationResponse response1 = new TokenValidationResponse(true);
        TokenValidationResponse response2 = new TokenValidationResponse(true);
        TokenValidationResponse response3 = new TokenValidationResponse(false);

        assertEquals(response1, response2);
        assertNotEquals(response1, response3);
        assertEquals(response1.hashCode(), response2.hashCode());
    }

    @Test
    void tokenValidationResponse_ShouldHaveProperToString() {
        TokenValidationResponse response = new TokenValidationResponse(true);

        String toString = response.toString();

        assertNotNull(toString);
        assertTrue(toString.contains("true"));
    }
}
