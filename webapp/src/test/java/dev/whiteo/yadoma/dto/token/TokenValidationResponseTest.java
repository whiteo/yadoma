package dev.whiteo.yadoma.dto.token;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class TokenValidationResponseTest {

    @Test
    void tokenValidationResponse_ShouldCreateWithValidTrue() {
        // Given
        Boolean valid = true;

        // When
        TokenValidationResponse response = new TokenValidationResponse(valid);

        // Then
        assertEquals(valid, response.valid());
        assertTrue(response.valid());
    }

    @Test
    void tokenValidationResponse_ShouldCreateWithValidFalse() {
        // Given
        Boolean valid = false;

        // When
        TokenValidationResponse response = new TokenValidationResponse(valid);

        // Then
        assertEquals(valid, response.valid());
        assertFalse(response.valid());
    }

    @Test
    void tokenValidationResponse_ShouldAllowNullValue() {
        // When
        TokenValidationResponse response = new TokenValidationResponse(null);

        // Then
        assertNull(response.valid());
    }

    @Test
    void tokenValidationResponse_ShouldSupportEquality() {
        // Given
        TokenValidationResponse response1 = new TokenValidationResponse(true);
        TokenValidationResponse response2 = new TokenValidationResponse(true);
        TokenValidationResponse response3 = new TokenValidationResponse(false);

        // Then
        assertEquals(response1, response2);
        assertNotEquals(response1, response3);
        assertEquals(response1.hashCode(), response2.hashCode());
    }

    @Test
    void tokenValidationResponse_ShouldHaveProperToString() {
        // Given
        TokenValidationResponse response = new TokenValidationResponse(true);

        // When
        String toString = response.toString();

        // Then
        assertNotNull(toString);
        assertTrue(toString.contains("true"));
    }
}
