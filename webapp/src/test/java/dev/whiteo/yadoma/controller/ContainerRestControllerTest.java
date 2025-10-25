package dev.whiteo.yadoma.controller;

import dev.whiteo.yadoma.dto.container.ContainerCreateRequest;
import dev.whiteo.yadoma.dto.container.ContainerResponse;
import dev.whiteo.yadoma.security.AuthInterceptor;
import dev.whiteo.yadoma.service.ContainerService;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

import java.time.LocalDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class ContainerRestControllerTest {

    @Mock
    private AuthInterceptor authInterceptor;

    @Mock
    private ContainerService containerService;

    @InjectMocks
    private ContainerRestController containerRestController;

    private static final String USER_ID = "user123";
    private static final String CONTAINER_ID = "container123";

    @BeforeEach
    void setUp() {
        when(authInterceptor.getUserId()).thenReturn(USER_ID);
    }

    @Test
    void findAll_ShouldReturnContainersList() {
        ContainerResponse container1 = new ContainerResponse("id1", "name1", LocalDateTime.now(), "running", "active");
        ContainerResponse container2 = new ContainerResponse("id2", "name2", LocalDateTime.now(), "stopped", "inactive");
        List<ContainerResponse> containers = Arrays.asList(container1, container2);

        when(containerService.findAll(anyString(), anyString())).thenReturn(containers);

        ResponseEntity<List<ContainerResponse>> response = containerRestController.findAll(USER_ID);

        assertEquals(HttpStatus.OK, response.getStatusCode());
        assertEquals(containers, response.getBody());
        verify(containerService).findAll(USER_ID, USER_ID);
    }

    @Test
    void getById_ShouldReturnContainer() {
        ContainerResponse container = new ContainerResponse("id1", "name1", LocalDateTime.now(), "running", "active");
        when(containerService.getById(anyString(), anyString())).thenReturn(container);

        ResponseEntity<ContainerResponse> response = containerRestController.getById(CONTAINER_ID);

        assertEquals(HttpStatus.OK, response.getStatusCode());
        assertEquals(container, response.getBody());
        verify(containerService).getById(CONTAINER_ID, USER_ID);
    }

    @Test
    void create_ShouldCreateContainer() {
        ContainerCreateRequest request = new ContainerCreateRequest("test-container", "nginx:latest", Collections.emptyList());
        doNothing().when(containerService).create(any(ContainerCreateRequest.class), anyString());

        ResponseEntity<Void> response = containerRestController.create(request);

        assertEquals(HttpStatus.OK, response.getStatusCode());
        verify(containerService).create(request, USER_ID);
    }

    @Test
    void delete_ShouldDeleteContainer() {
        doNothing().when(containerService).delete(anyString(), anyString());

        ResponseEntity<Void> response = containerRestController.delete(CONTAINER_ID);

        assertEquals(HttpStatus.NO_CONTENT, response.getStatusCode());
        verify(containerService).delete(CONTAINER_ID, USER_ID);
    }

    @Test
    void start_ShouldStartContainer() {
        doNothing().when(containerService).start(anyString(), anyString());

        ResponseEntity<Void> response = containerRestController.start(CONTAINER_ID);

        assertEquals(HttpStatus.OK, response.getStatusCode());
        verify(containerService).start(CONTAINER_ID, USER_ID);
    }

    @Test
    void stop_ShouldStopContainer() {
        doNothing().when(containerService).stop(anyString(), anyString());

        ResponseEntity<Void> response = containerRestController.stop(CONTAINER_ID);

        assertEquals(HttpStatus.OK, response.getStatusCode());
        verify(containerService).stop(CONTAINER_ID, USER_ID);
    }

    @Test
    void restart_ShouldRestartContainer() {
        doNothing().when(containerService).restart(anyString(), anyString());

        ResponseEntity<Void> response = containerRestController.restart(CONTAINER_ID);

        assertEquals(HttpStatus.OK, response.getStatusCode());
        verify(containerService).restart(CONTAINER_ID, USER_ID);
    }

    @Test
    void getLogs_ShouldReturnStreamingResponseBody() throws Exception {
        container.v1.Container.GetContainerLogsResponse logResponse1 = container.v1.Container.GetContainerLogsResponse.newBuilder()
                .setChunk(com.google.protobuf.ByteString.copyFromUtf8("Log line 1\n"))
                .build();
        container.v1.Container.GetContainerLogsResponse logResponse2 = container.v1.Container.GetContainerLogsResponse.newBuilder()
                .setChunk(com.google.protobuf.ByteString.copyFromUtf8("Log line 2\n"))
                .build();

        java.util.Iterator<container.v1.Container.GetContainerLogsResponse> mockIterator = Arrays.asList(logResponse1, logResponse2).iterator();
        when(containerService.getLogs(CONTAINER_ID, USER_ID, false)).thenReturn(mockIterator);

        org.springframework.http.ResponseEntity<org.springframework.web.servlet.mvc.method.annotation.StreamingResponseBody> response =
                containerRestController.getLogs(CONTAINER_ID, false);

        assertEquals(HttpStatus.OK, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(org.springframework.http.MediaType.TEXT_PLAIN, response.getHeaders().getContentType());

        java.io.ByteArrayOutputStream outputStream = new java.io.ByteArrayOutputStream();
        response.getBody().writeTo(outputStream);
        String result = outputStream.toString(java.nio.charset.StandardCharsets.UTF_8);
        assertEquals("Log line 1\nLog line 2\n", result);

        verify(containerService).getLogs(CONTAINER_ID, USER_ID, false);
    }

    @Test
    void getLogs_ShouldReturnStreamingResponseBodyWithFollow() throws Exception {
        container.v1.Container.GetContainerLogsResponse logResponse = container.v1.Container.GetContainerLogsResponse.newBuilder()
                .setChunk(com.google.protobuf.ByteString.copyFromUtf8("Following logs\n"))
                .build();

        java.util.Iterator<container.v1.Container.GetContainerLogsResponse> mockIterator = Collections.singletonList(logResponse).iterator();
        when(containerService.getLogs(CONTAINER_ID, USER_ID, true)).thenReturn(mockIterator);

        org.springframework.http.ResponseEntity<org.springframework.web.servlet.mvc.method.annotation.StreamingResponseBody> response =
                containerRestController.getLogs(CONTAINER_ID, true);

        assertEquals(HttpStatus.OK, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(org.springframework.http.MediaType.TEXT_PLAIN, response.getHeaders().getContentType());

        java.io.ByteArrayOutputStream outputStream = new java.io.ByteArrayOutputStream();
        response.getBody().writeTo(outputStream);
        String result = outputStream.toString(java.nio.charset.StandardCharsets.UTF_8);
        assertEquals("Following logs\n", result);

        verify(containerService).getLogs(CONTAINER_ID, USER_ID, true);
    }

    @Test
    void getLogs_ShouldHandleExceptionDuringStreaming() throws Exception {
        java.util.Iterator<container.v1.Container.GetContainerLogsResponse> mockIterator = mock(java.util.Iterator.class);
        when(mockIterator.hasNext()).thenReturn(true);
        when(mockIterator.next()).thenThrow(new RuntimeException("Log streaming error"));
        when(containerService.getLogs(CONTAINER_ID, USER_ID, false)).thenReturn(mockIterator);

        org.springframework.http.ResponseEntity<org.springframework.web.servlet.mvc.method.annotation.StreamingResponseBody> response =
                containerRestController.getLogs(CONTAINER_ID, false);

        assertEquals(HttpStatus.OK, response.getStatusCode());
        assertNotNull(response.getBody());

        java.io.ByteArrayOutputStream outputStream = new java.io.ByteArrayOutputStream();
        response.getBody().writeTo(outputStream);
        String result = outputStream.toString(java.nio.charset.StandardCharsets.UTF_8);
        assertTrue(result.contains("Error retrieving logs"));
        assertTrue(result.contains("Log streaming error"));
    }

    @Test
    void getLogs_ShouldHandleEmptyLogStream() throws Exception {
        java.util.Iterator<container.v1.Container.GetContainerLogsResponse> mockIterator = Collections.emptyIterator();
        when(containerService.getLogs(CONTAINER_ID, USER_ID, false)).thenReturn(mockIterator);

        org.springframework.http.ResponseEntity<org.springframework.web.servlet.mvc.method.annotation.StreamingResponseBody> response =
                containerRestController.getLogs(CONTAINER_ID, false);

        assertEquals(HttpStatus.OK, response.getStatusCode());
        assertNotNull(response.getBody());

        java.io.ByteArrayOutputStream outputStream = new java.io.ByteArrayOutputStream();
        response.getBody().writeTo(outputStream);
        String result = outputStream.toString(java.nio.charset.StandardCharsets.UTF_8);
        assertEquals("", result);
    }
}
