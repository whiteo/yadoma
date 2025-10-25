package dev.whiteo.yadoma.websocket;

import container.v1.Container;
import container.v1.ContainerServiceGrpc;
import dev.whiteo.yadoma.security.TokenInteract;
import dev.whiteo.yadoma.service.UserService;
import io.grpc.stub.StreamObserver;
import jakarta.websocket.CloseReason;
import jakarta.websocket.RemoteEndpoint;
import jakarta.websocket.Session;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.ArgumentCaptor;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.io.IOException;
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

    @Mock
    private RemoteEndpoint.Basic basicRemote;

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
        lenient().when(session.getBasicRemote()).thenReturn(basicRemote);
    }

    @Test
    void onOpen_ShouldValidateTokenAndStartStreaming() throws Exception {
        when(tokenInteract.validateToken(TOKEN)).thenReturn(true);
        when(tokenInteract.getUserId(TOKEN)).thenReturn(USER_ID);
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);

        assertDoesNotThrow(() -> containerStatsEndpoint.onOpen(session, CONTAINER_ID));

        verify(tokenInteract).validateToken(TOKEN);
        verify(tokenInteract).getUserId(TOKEN);
        verify(userService).isContainerIdContains(CONTAINER_ID, USER_ID);
    }

    @Test
    void onOpen_ShouldCloseSessionOnInvalidToken() throws Exception {
        when(tokenInteract.validateToken(TOKEN)).thenReturn(false);

        containerStatsEndpoint.onOpen(session, CONTAINER_ID);

        verify(tokenInteract).validateToken(TOKEN);
        verify(session).close(any(CloseReason.class));
        verify(userService, never()).isContainerIdContains(any(), any());
    }

    @Test
    void onOpen_ShouldCloseSessionOnUnauthorizedContainer() throws Exception {
        when(tokenInteract.validateToken(TOKEN)).thenReturn(true);
        when(tokenInteract.getUserId(TOKEN)).thenReturn(USER_ID);
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(false);

        containerStatsEndpoint.onOpen(session, CONTAINER_ID);

        verify(tokenInteract).validateToken(TOKEN);
        verify(tokenInteract).getUserId(TOKEN);
        verify(userService).isContainerIdContains(CONTAINER_ID, USER_ID);
        verify(session).close(any(CloseReason.class));
    }

    @Test
    void onClose_ShouldExecuteWithoutErrors() {
        assertDoesNotThrow(() -> containerStatsEndpoint.onClose(session));
    }

    @Test
    void onOpen_ShouldHandleNullToken() throws Exception {
        Map<String, List<String>> requestParameterMap = new HashMap<>();
        requestParameterMap.put("token", null);
        when(session.getRequestParameterMap()).thenReturn(requestParameterMap);

        containerStatsEndpoint.onOpen(session, CONTAINER_ID);

        verify(session).close(any(CloseReason.class));
        verifyNoInteractions(tokenInteract);
    }

    @Test
    void onOpen_ShouldHandleEmptyTokenList() throws Exception {
        Map<String, List<String>> requestParameterMap = new HashMap<>();
        requestParameterMap.put("token", Collections.emptyList());
        when(session.getRequestParameterMap()).thenReturn(requestParameterMap);

        containerStatsEndpoint.onOpen(session, CONTAINER_ID);

        verify(session).close(any(CloseReason.class));
        verifyNoInteractions(tokenInteract);
    }

    @Test
    void onOpen_ShouldHandleTokenValidationException() throws Exception {
        when(tokenInteract.validateToken(TOKEN)).thenThrow(new RuntimeException("Token validation error"));

        containerStatsEndpoint.onOpen(session, CONTAINER_ID);

        verify(tokenInteract).validateToken(TOKEN);
    }

    @Test
    void streamObserver_OnNext_ShouldSendFormattedStatsAsJson() throws Exception {
        ArgumentCaptor<StreamObserver<Container.GetContainerStatsResponse>> observerCaptor =
                ArgumentCaptor.forClass(StreamObserver.class);

        when(tokenInteract.validateToken(TOKEN)).thenReturn(true);
        when(tokenInteract.getUserId(TOKEN)).thenReturn(USER_ID);
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);
        doAnswer(invocation -> {
            return null;
        }).when(containerStub).getContainerStats(any(Container.GetContainerStatsRequest.class), observerCaptor.capture());

        containerStatsEndpoint.onOpen(session, CONTAINER_ID);
        Thread.sleep(100);

        StreamObserver<Container.GetContainerStatsResponse> observer = observerCaptor.getValue();
        assertNotNull(observer);

        Container.GetContainerStatsResponse statsResponse = Container.GetContainerStatsResponse.newBuilder()
                .setCpuUsage(75)
                .setMemUsage(1024000000)
                .setMemLimit(2048000000)
                .setNetInput(500000)
                .setNetOutput(300000)
                .build();

        observer.onNext(statsResponse);
        Thread.sleep(50);

        ArgumentCaptor<String> jsonCaptor = ArgumentCaptor.forClass(String.class);
        verify(basicRemote, atLeastOnce()).sendText(jsonCaptor.capture());

        String sentJson = jsonCaptor.getValue();
        assertTrue(sentJson.contains("\"cpu\": 75"));
        assertTrue(sentJson.contains("\"memUsage\": 1024000000"));
        assertTrue(sentJson.contains("\"memLimit\": 2048000000"));
        assertTrue(sentJson.contains("\"netInput\": 500000"));
        assertTrue(sentJson.contains("\"netOutput\": 300000"));
    }

    @Test
    void streamObserver_OnNext_ShouldHandleIOException() throws Exception {
        ArgumentCaptor<StreamObserver<Container.GetContainerStatsResponse>> observerCaptor =
                ArgumentCaptor.forClass(StreamObserver.class);

        when(tokenInteract.validateToken(TOKEN)).thenReturn(true);
        when(tokenInteract.getUserId(TOKEN)).thenReturn(USER_ID);
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);
        doAnswer(invocation -> {
            return null;
        }).when(containerStub).getContainerStats(any(Container.GetContainerStatsRequest.class), observerCaptor.capture());
        doThrow(new IOException("Connection lost")).when(basicRemote).sendText(anyString());

        containerStatsEndpoint.onOpen(session, CONTAINER_ID);
        Thread.sleep(100);

        StreamObserver<Container.GetContainerStatsResponse> observer = observerCaptor.getValue();
        assertNotNull(observer);

        Container.GetContainerStatsResponse statsResponse = Container.GetContainerStatsResponse.newBuilder()
                .setCpuUsage(50)
                .setMemUsage(512000000)
                .setMemLimit(1024000000)
                .setNetInput(100000)
                .setNetOutput(50000)
                .build();

        observer.onNext(statsResponse);
        Thread.sleep(50);

        verify(session, atLeastOnce()).close();
    }

    @Test
    void streamObserver_OnError_ShouldSendErrorMessageAndCloseSession() throws Exception {
        ArgumentCaptor<StreamObserver<Container.GetContainerStatsResponse>> observerCaptor =
                ArgumentCaptor.forClass(StreamObserver.class);

        when(tokenInteract.validateToken(TOKEN)).thenReturn(true);
        when(tokenInteract.getUserId(TOKEN)).thenReturn(USER_ID);
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);
        doAnswer(invocation -> {
            return null;
        }).when(containerStub).getContainerStats(any(Container.GetContainerStatsRequest.class), observerCaptor.capture());

        containerStatsEndpoint.onOpen(session, CONTAINER_ID);
        Thread.sleep(100);

        StreamObserver<Container.GetContainerStatsResponse> observer = observerCaptor.getValue();
        assertNotNull(observer);

        RuntimeException testError = new RuntimeException("gRPC stream error");
        observer.onError(testError);
        Thread.sleep(50);

        ArgumentCaptor<String> messageCaptor = ArgumentCaptor.forClass(String.class);
        verify(basicRemote, atLeastOnce()).sendText(messageCaptor.capture());
        String errorMessage = messageCaptor.getValue();
        assertTrue(errorMessage.contains("[error]"));
        assertTrue(errorMessage.contains("gRPC stream error"));

        verify(session, atLeastOnce()).close();
    }

    @Test
    void streamObserver_OnCompleted_ShouldCloseSession() throws Exception {
        ArgumentCaptor<StreamObserver<Container.GetContainerStatsResponse>> observerCaptor =
                ArgumentCaptor.forClass(StreamObserver.class);

        when(tokenInteract.validateToken(TOKEN)).thenReturn(true);
        when(tokenInteract.getUserId(TOKEN)).thenReturn(USER_ID);
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);
        doAnswer(invocation -> {
            return null;
        }).when(containerStub).getContainerStats(any(Container.GetContainerStatsRequest.class), observerCaptor.capture());

        containerStatsEndpoint.onOpen(session, CONTAINER_ID);
        Thread.sleep(100);

        StreamObserver<Container.GetContainerStatsResponse> observer = observerCaptor.getValue();
        assertNotNull(observer);

        observer.onCompleted();
        Thread.sleep(50);

        verify(session, atLeastOnce()).close();
    }

    @Test
    void streamObserver_OnError_ShouldHandleIOExceptionWhenSendingError() throws Exception {
        ArgumentCaptor<StreamObserver<Container.GetContainerStatsResponse>> observerCaptor =
                ArgumentCaptor.forClass(StreamObserver.class);

        when(tokenInteract.validateToken(TOKEN)).thenReturn(true);
        when(tokenInteract.getUserId(TOKEN)).thenReturn(USER_ID);
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);
        doAnswer(invocation -> {
            return null;
        }).when(containerStub).getContainerStats(any(Container.GetContainerStatsRequest.class), observerCaptor.capture());
        doThrow(new IOException("Cannot send error")).when(basicRemote).sendText(anyString());

        containerStatsEndpoint.onOpen(session, CONTAINER_ID);
        Thread.sleep(100);

        StreamObserver<Container.GetContainerStatsResponse> observer = observerCaptor.getValue();
        assertNotNull(observer);

        RuntimeException testError = new RuntimeException("gRPC stream error");
        observer.onError(testError);
        Thread.sleep(50);

        verify(session, atLeastOnce()).close();
    }
}
