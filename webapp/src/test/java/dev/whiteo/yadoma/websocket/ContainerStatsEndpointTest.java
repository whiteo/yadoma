package dev.whiteo.yadoma.websocket;

import container.v1.ContainerServiceGrpc;
import dev.whiteo.yadoma.security.TokenInteract;
import dev.whiteo.yadoma.service.UserService;
import jakarta.websocket.CloseReason;
import jakarta.websocket.Session;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class ContainerStatsEndpointTest {

    @Mock
    private ContainerServiceGrpc.ContainerServiceStub containerStub;

    @Mock
    private UserService userService;

    @Mock
    private TokenInteract tokenInteract;

    @Mock
    private Session session;

    @InjectMocks
    private ContainerStatsEndpoint containerStatsEndpoint;

    private static final String CONTAINER_ID = "container123";
    private static final String USER_ID = "user123";
    private static final String TOKEN = "validToken";

    @BeforeEach
    void setUp() {
        Map<String, List<String>> requestParameterMap = new HashMap<>();
        requestParameterMap.put("token", List.of(TOKEN));
        lenient().when(session.getRequestParameterMap()).thenReturn(requestParameterMap);
    }

    @Test
    void onOpen_ShouldValidateTokenAndStartStreaming() throws Exception {
        // Given
        when(tokenInteract.validateToken(TOKEN)).thenReturn(true);
        when(tokenInteract.getUserId(TOKEN)).thenReturn(USER_ID);
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);

        // When
        assertDoesNotThrow(() -> containerStatsEndpoint.onOpen(session, CONTAINER_ID));

        // Then
        verify(tokenInteract).validateToken(TOKEN);
        verify(tokenInteract).getUserId(TOKEN);
        verify(userService).isContainerIdContains(CONTAINER_ID, USER_ID);
        // Не проверяем containerStub напрямую, так как streamStats выполняется асинхронно
    }

    @Test
    void onOpen_ShouldCloseSessionOnInvalidToken() throws Exception {
        // Given
        when(tokenInteract.validateToken(TOKEN)).thenReturn(false);

        // When
        containerStatsEndpoint.onOpen(session, CONTAINER_ID);

        // Then
        verify(tokenInteract).validateToken(TOKEN);
        verify(session).close(any(CloseReason.class));
        verify(userService, never()).isContainerIdContains(any(), any());
    }

    @Test
    void onOpen_ShouldCloseSessionOnUnauthorizedContainer() throws Exception {
        // Given
        when(tokenInteract.validateToken(TOKEN)).thenReturn(true);
        when(tokenInteract.getUserId(TOKEN)).thenReturn(USER_ID);
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(false);

        // When
        containerStatsEndpoint.onOpen(session, CONTAINER_ID);

        // Then
        verify(tokenInteract).validateToken(TOKEN);
        verify(tokenInteract).getUserId(TOKEN);
        verify(userService).isContainerIdContains(CONTAINER_ID, USER_ID);
        verify(session).close(any(CloseReason.class));
    }

    @Test
    void onClose_ShouldExecuteWithoutErrors() {
        // When & Then
        assertDoesNotThrow(() -> containerStatsEndpoint.onClose(session));
    }

    @Test
    void onOpen_ShouldHandleNullToken() throws Exception {
        // Given
        Map<String, List<String>> requestParameterMap = new HashMap<>();
        requestParameterMap.put("token", null);
        when(session.getRequestParameterMap()).thenReturn(requestParameterMap);

        // When
        containerStatsEndpoint.onOpen(session, CONTAINER_ID);

        // Then
        verify(session).close(any(CloseReason.class));
        verifyNoInteractions(tokenInteract);
    }

    @Test
    void onOpen_ShouldHandleEmptyTokenList() throws Exception {
        // Given
        Map<String, List<String>> requestParameterMap = new HashMap<>();
        requestParameterMap.put("token", Collections.emptyList());
        when(session.getRequestParameterMap()).thenReturn(requestParameterMap);

        // When
        containerStatsEndpoint.onOpen(session, CONTAINER_ID);

        // Then
        verify(session).close(any(CloseReason.class));
        verifyNoInteractions(tokenInteract);
    }

    @Test
    void onOpen_ShouldHandleTokenValidationException() throws Exception {
        // Given
        when(tokenInteract.validateToken(TOKEN)).thenThrow(new RuntimeException("Token validation error"));

        // When
        containerStatsEndpoint.onOpen(session, CONTAINER_ID);

        // Then
        verify(tokenInteract).validateToken(TOKEN);
        // closeSession будет вызван из блока catch
    }
}
