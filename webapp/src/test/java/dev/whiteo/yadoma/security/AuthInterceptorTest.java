package dev.whiteo.yadoma.security;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class AuthInterceptorTest {

    @Mock
    private TokenInteract tokenInteract;

    @Mock
    private HttpServletRequest request;

    @Mock
    private HttpServletResponse response;

    @InjectMocks
    private AuthInterceptor authInterceptor;

    private Object handler = new Object();

    @BeforeEach
    void setUp() {}

    @Test
    void preHandle_WithVersionEndpoint_ShouldReturnTrue() throws Exception {
        when(request.getRequestURI()).thenReturn("/version");

        boolean result = authInterceptor.preHandle(request, response, handler);

        assertTrue(result);
        verifyNoInteractions(tokenInteract);
    }

    @Test
    void preHandle_WithAuthenticateEndpoint_ShouldReturnTrue() throws Exception {
        when(request.getRequestURI()).thenReturn("/api/v1/authenticate");

        boolean result = authInterceptor.preHandle(request, response, handler);

        assertTrue(result);
        verifyNoInteractions(tokenInteract);
    }

    @Test
    void preHandle_WithUserCreateEndpoint_ShouldReturnTrue() throws Exception {
        when(request.getRequestURI()).thenReturn("/api/v1/user/create");

        boolean result = authInterceptor.preHandle(request, response, handler);

        assertTrue(result);
        verifyNoInteractions(tokenInteract);
    }

    @Test
    void preHandle_WithValidToken_ShouldReturnTrue() throws Exception {
        String token = "validToken";
        String userId = "user123";

        when(request.getRequestURI()).thenReturn("/api/v1/containers");
        when(tokenInteract.getToken(request)).thenReturn(token);
        when(tokenInteract.validateToken(token)).thenReturn(true);
        when(tokenInteract.getUserId(token)).thenReturn(userId);

        boolean result = authInterceptor.preHandle(request, response, handler);

        assertTrue(result);
        assertEquals(userId, authInterceptor.getUserId());
        verify(tokenInteract).getToken(request);
        verify(tokenInteract).validateToken(token);
        verify(tokenInteract).getUserId(token);
    }

    @Test
    void preHandle_WithInvalidToken_ShouldReturnFalse() throws Exception {
        String token = "invalidToken";

        when(request.getRequestURI()).thenReturn("/api/v1/containers");
        when(tokenInteract.getToken(request)).thenReturn(token);
        when(tokenInteract.validateToken(token)).thenReturn(false);

        boolean result = authInterceptor.preHandle(request, response, handler);

        assertFalse(result);
        verify(response).sendError(HttpServletResponse.SC_UNAUTHORIZED, "Unauthorized");
        verify(tokenInteract).getToken(request);
        verify(tokenInteract).validateToken(token);
        verify(tokenInteract, never()).getUserId(token);
    }

    @Test
    void preHandle_WithTokenButEmptyUserId_ShouldReturnFalse() throws Exception {
        String token = "validToken";

        when(request.getRequestURI()).thenReturn("/api/v1/containers");
        when(tokenInteract.getToken(request)).thenReturn(token);
        when(tokenInteract.validateToken(token)).thenReturn(true);
        when(tokenInteract.getUserId(token)).thenReturn("");

        boolean result = authInterceptor.preHandle(request, response, handler);

        assertFalse(result);
        verify(response).sendError(HttpServletResponse.SC_UNAUTHORIZED, "Unauthorized");
    }

    @Test
    void preHandle_WithTokenButNullUserId_ShouldReturnFalse() throws Exception {
        String token = "validToken";

        when(request.getRequestURI()).thenReturn("/api/v1/containers");
        when(tokenInteract.getToken(request)).thenReturn(token);
        when(tokenInteract.validateToken(token)).thenReturn(true);
        when(tokenInteract.getUserId(token)).thenReturn(null);

        boolean result = authInterceptor.preHandle(request, response, handler);

        assertFalse(result);
        verify(response).sendError(HttpServletResponse.SC_UNAUTHORIZED, "Unauthorized");
    }

    @Test
    void postHandle_ShouldClearThreadLocal() {
        authInterceptor.postHandle(request, response, handler, null);

        assertDoesNotThrow(() -> authInterceptor.getUserId());
    }

    @Test
    void getUserId_ShouldReturnCurrentUserId() {
        assertDoesNotThrow(() -> authInterceptor.getUserId());
    }
}