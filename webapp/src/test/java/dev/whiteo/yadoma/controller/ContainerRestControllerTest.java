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
        // Given
        ContainerResponse container1 = new ContainerResponse("id1", "name1", LocalDateTime.now(), "running", "active");
        ContainerResponse container2 = new ContainerResponse("id2", "name2", LocalDateTime.now(), "stopped", "inactive");
        List<ContainerResponse> containers = Arrays.asList(container1, container2);

        when(containerService.findAll(anyString(), anyString())).thenReturn(containers);

        // When
        ResponseEntity<List<ContainerResponse>> response = containerRestController.findAll(USER_ID);

        // Then
        assertEquals(HttpStatus.OK, response.getStatusCode());
        assertEquals(containers, response.getBody());
        verify(containerService).findAll(USER_ID, USER_ID);
    }

    @Test
    void getById_ShouldReturnContainer() {
        // Given
        ContainerResponse container = new ContainerResponse("id1", "name1", LocalDateTime.now(), "running", "active");
        when(containerService.getById(anyString(), anyString())).thenReturn(container);

        // When
        ResponseEntity<ContainerResponse> response = containerRestController.getById(CONTAINER_ID);

        // Then
        assertEquals(HttpStatus.OK, response.getStatusCode());
        assertEquals(container, response.getBody());
        verify(containerService).getById(CONTAINER_ID, USER_ID);
    }

    @Test
    void create_ShouldCreateContainer() {
        // Given
        ContainerCreateRequest request = new ContainerCreateRequest("test-container", "nginx:latest", Collections.emptyList());
        doNothing().when(containerService).create(any(ContainerCreateRequest.class), anyString());

        // When
        ResponseEntity<Void> response = containerRestController.create(request);

        // Then
        assertEquals(HttpStatus.OK, response.getStatusCode());
        verify(containerService).create(request, USER_ID);
    }

    @Test
    void delete_ShouldDeleteContainer() {
        // Given
        doNothing().when(containerService).delete(anyString(), anyString());

        // When
        ResponseEntity<Void> response = containerRestController.delete(CONTAINER_ID);

        // Then
        assertEquals(HttpStatus.NO_CONTENT, response.getStatusCode());
        verify(containerService).delete(CONTAINER_ID, USER_ID);
    }

    @Test
    void start_ShouldStartContainer() {
        // Given
        doNothing().when(containerService).start(anyString(), anyString());

        // When
        ResponseEntity<Void> response = containerRestController.start(CONTAINER_ID);

        // Then
        assertEquals(HttpStatus.OK, response.getStatusCode());
        verify(containerService).start(CONTAINER_ID, USER_ID);
    }

    @Test
    void stop_ShouldStopContainer() {
        // Given
        doNothing().when(containerService).stop(anyString(), anyString());

        // When
        ResponseEntity<Void> response = containerRestController.stop(CONTAINER_ID);

        // Then
        assertEquals(HttpStatus.OK, response.getStatusCode());
        verify(containerService).stop(CONTAINER_ID, USER_ID);
    }

    @Test
    void restart_ShouldRestartContainer() {
        // Given
        doNothing().when(containerService).restart(anyString(), anyString());

        // When
        ResponseEntity<Void> response = containerRestController.restart(CONTAINER_ID);

        // Then
        assertEquals(HttpStatus.OK, response.getStatusCode());
        verify(containerService).restart(CONTAINER_ID, USER_ID);
    }
}
