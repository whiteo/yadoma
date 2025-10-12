package dev.whiteo.yadoma.security;

import dev.whiteo.yadoma.exception.ExecutionConflictException;
import jakarta.servlet.http.HttpServletRequest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.http.HttpHeaders;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.test.util.ReflectionTestUtils;
import org.springframework.web.client.HttpClientErrorException;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class TokenInteractTest {

    @InjectMocks
    private TokenInteract tokenInteract;

    @Mock
    private HttpServletRequest request;

    @Mock
    private UserDetails userDetails;

    private final String secretKey = "mySecretKeyForTestingPurposesOnlyAndShouldBeAtLeast256Bits";
    private final long expirationTime = 86400000L;

    @BeforeEach
    void setUp() {
        ReflectionTestUtils.setField(tokenInteract, "secretKey", secretKey);
        ReflectionTestUtils.setField(tokenInteract, "expirationTime", expirationTime);
    }

    @Test
    void generateToken_ShouldCreateValidToken() {
        when(userDetails.getUsername()).thenReturn("testUser");

        String token = tokenInteract.generateToken(userDetails);

        assertNotNull(token);
        assertFalse(token.isEmpty());
        assertEquals(3, token.split("\\.").length);
    }

    @Test
    void validateToken_WithValidToken_ShouldReturnTrue() {
        when(userDetails.getUsername()).thenReturn("testUser");
        String token = tokenInteract.generateToken(userDetails);

        boolean isValid = tokenInteract.validateToken(token);

        assertTrue(isValid);
    }

    @Test
    void validateToken_WithInvalidToken_ShouldThrowException() {
        String invalidToken = "invalid.token.here";

        assertThrows(ExecutionConflictException.class, () ->
            tokenInteract.validateToken(invalidToken));
    }

    @Test
    void validateToken_WithMalformedToken_ShouldThrowException() {
        String malformedToken = "notAValidJwtToken";

        assertThrows(ExecutionConflictException.class, () ->
            tokenInteract.validateToken(malformedToken));
    }

    @Test
    void validateToken_WithEmptyToken_ShouldThrowException() {
        String emptyToken = "";

        assertThrows(ExecutionConflictException.class, () ->
            tokenInteract.validateToken(emptyToken));
    }

    @Test
    void getUserId_WithValidToken_ShouldReturnUserId() {
        String expectedUserId = "testUser123";
        when(userDetails.getUsername()).thenReturn(expectedUserId);
        String token = tokenInteract.generateToken(userDetails);

        String actualUserId = tokenInteract.getUserId(token);

        assertEquals(expectedUserId, actualUserId);
    }

    @Test
    void getToken_WithValidBearerToken_ShouldReturnToken() {
        String token = "validTokenString";
        String authHeader = "Bearer " + token;
        when(request.getHeader(HttpHeaders.AUTHORIZATION)).thenReturn(authHeader);

        String extractedToken = tokenInteract.getToken(request);

        assertEquals(token, extractedToken);
    }

    @Test
    void getToken_WithoutBearerPrefix_ShouldThrowException() {
        String authHeader = "InvalidTokenFormat";
        when(request.getHeader(HttpHeaders.AUTHORIZATION)).thenReturn(authHeader);

        assertThrows(HttpClientErrorException.class, () ->
            tokenInteract.getToken(request));
    }

    @Test
    void getToken_WithNullHeader_ShouldThrowException() {
        when(request.getHeader(HttpHeaders.AUTHORIZATION)).thenReturn(null);

        assertThrows(HttpClientErrorException.class, () ->
            tokenInteract.getToken(request));
    }

    @Test
    void getToken_WithEmptyHeader_ShouldThrowException() {
        when(request.getHeader(HttpHeaders.AUTHORIZATION)).thenReturn("");

        assertThrows(HttpClientErrorException.class, () ->
            tokenInteract.getToken(request));
    }

    @Test
    void generateToken_WithDifferentUsers_ShouldProduceDifferentTokens() {
        UserDetails user1 = org.mockito.Mockito.mock(UserDetails.class);
        UserDetails user2 = org.mockito.Mockito.mock(UserDetails.class);
        when(user1.getUsername()).thenReturn("user1");
        when(user2.getUsername()).thenReturn("user2");

        String token1 = tokenInteract.generateToken(user1);
        String token2 = tokenInteract.generateToken(user2);

        assertNotEquals(token1, token2);
    }
}