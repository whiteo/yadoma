package dev.whiteo.yadoma.service;

import container.v1.Container;
import container.v1.ContainerServiceGrpc;
import dev.whiteo.yadoma.domain.User;
import dev.whiteo.yadoma.dto.container.ContainerCreateRequest;
import dev.whiteo.yadoma.dto.container.ContainerResponse;
import dev.whiteo.yadoma.mapper.ContainerMapper;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.security.authentication.BadCredentialsException;

import java.time.LocalDateTime;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.verifyNoInteractions;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class ContainerServiceTest {

    @Mock
    private ContainerServiceGrpc.ContainerServiceBlockingStub containerStub;

    @Mock
    private image.v1.ImageServiceGrpc.ImageServiceBlockingStub imageStub;

    @Mock
    private UserService userService;

    @Mock
    private ContainerMapper mapper;

    @InjectMocks
    private ContainerService containerService;

    private static final String USER_ID = "user123";
    private static final String ADMIN_ID = "admin123";
    private static final String CONTAINER_ID = "container123";

    private User testUser;
    private Container.GetContainersResponse getContainersResponse;

    @BeforeEach
    void setUp() {
        testUser = new User();
        testUser.setId(USER_ID);
        testUser.setContainerIds(List.of(CONTAINER_ID));

        Container.GetContainerResponse containerResponse = Container.GetContainerResponse.newBuilder()
                .setId(CONTAINER_ID)
                .addNames("test-container")
                .setImage("nginx:latest")
                .setState("running")
                .build();

        getContainersResponse = Container.GetContainersResponse.newBuilder()
                .addContainers(containerResponse)
                .build();
    }

    @Test
    void findAll_ShouldReturnFilteredContainers() {
        when(userService.validateUserAccess(USER_ID, ADMIN_ID)).thenReturn(testUser);
        when(containerStub.getContainers(any())).thenReturn(getContainersResponse);

        List<ContainerResponse> result = containerService.findAll(USER_ID, ADMIN_ID);

        verify(userService).validateUserAccess(USER_ID, ADMIN_ID);
        verify(containerStub).getContainers(any());
    }

    @Test
    void findAll_ShouldReturnNullOnGrpcError() {
        when(userService.validateUserAccess(USER_ID, ADMIN_ID)).thenReturn(testUser);
        when(containerStub.getContainers(any())).thenThrow(new RuntimeException("gRPC error"));

        List<ContainerResponse> result = containerService.findAll(USER_ID, ADMIN_ID);

        assertNull(result);
    }

    @Test
    void getById_ShouldReturnContainer() {
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);
        Container.GetContainerDetailsResponse getContainerResponse = Container.GetContainerDetailsResponse.newBuilder()
                .setId(CONTAINER_ID)
                .setName("test-container")
                .setImage("nginx:latest")
                .setState("running")
                .build();
        when(containerStub.getContainerDetails(any())).thenReturn(getContainerResponse);
        ContainerResponse containerResponse = new ContainerResponse("id1", "name1", LocalDateTime.now(), "running", "active");
        when(mapper.toResponseDTO(any(Container.GetContainerDetailsResponse.class))).thenReturn(containerResponse);

        ContainerResponse result = containerService.getById(CONTAINER_ID, USER_ID);

        assertNotNull(result);
        verify(userService).isContainerIdContains(CONTAINER_ID, USER_ID);
        verify(containerStub).getContainerDetails(any());
    }

    @Test
    void getById_ShouldThrowBadCredentialsWhenUserNotAuthorized() {
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(false);

        assertThrows(BadCredentialsException.class, () ->
            containerService.getById(CONTAINER_ID, USER_ID));

        verify(userService).isContainerIdContains(CONTAINER_ID, USER_ID);
        verifyNoInteractions(containerStub);
    }

    @Test
    void create_ShouldCreateContainer() {
        ContainerCreateRequest request = new ContainerCreateRequest("test-container", "nginx:latest", Collections.emptyList());
        when(userService.getUserById(USER_ID)).thenReturn(testUser);

        image.v1.Image.GetImagesResponse imagesResponse = image.v1.Image.GetImagesResponse.newBuilder()
                .addImages(image.v1.Image.GetImageResponse.newBuilder()
                        .addRepoTags("nginx:latest")
                        .build())
                .build();
        when(imageStub.getImages(any())).thenReturn(imagesResponse);

        Container.CreateContainerResponse createResponse = Container.CreateContainerResponse.newBuilder()
                .setId(CONTAINER_ID)
                .build();
        when(containerStub.createContainer(any())).thenReturn(createResponse);

        assertDoesNotThrow(() -> containerService.create(request, USER_ID));

        verify(userService).getUserById(USER_ID);
        verify(imageStub).getImages(any());
        verify(containerStub).createContainer(any());
        verify(userService).addContainerToUser(CONTAINER_ID, testUser);
    }

    @Test
    void create_ShouldHandleGrpcErrorGracefully() {
        ContainerCreateRequest request = new ContainerCreateRequest("test-container", "nginx:latest", Collections.emptyList());
        when(userService.getUserById(USER_ID)).thenReturn(testUser);

        image.v1.Image.GetImagesResponse imagesResponse = image.v1.Image.GetImagesResponse.newBuilder()
                .addImages(image.v1.Image.GetImageResponse.newBuilder()
                        .addRepoTags("nginx:latest")
                        .build())
                .build();
        when(imageStub.getImages(any())).thenReturn(imagesResponse);

        when(containerStub.createContainer(any())).thenThrow(new io.grpc.StatusRuntimeException(io.grpc.Status.INTERNAL));

        assertThrows(dev.whiteo.yadoma.exception.ExecutionConflictException.class,
                () -> containerService.create(request, USER_ID));
    }

    @Test
    void delete_ShouldDeleteContainer() {
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);
        when(userService.getUserById(USER_ID)).thenReturn(testUser);
        Container.RemoveContainerResponse removeResponse = Container.RemoveContainerResponse.newBuilder()
                .setSuccess(true)
                .build();
        when(containerStub.removeContainer(any())).thenReturn(removeResponse);

        assertDoesNotThrow(() -> containerService.delete(CONTAINER_ID, USER_ID));

        verify(userService).getUserById(USER_ID);
        verify(containerStub).removeContainer(any());
        verify(userService).removeContainerFromUser(CONTAINER_ID, testUser);
    }

    @Test
    void delete_ShouldThrowWhenRemoveFails() {
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);
        when(userService.getUserById(USER_ID)).thenReturn(testUser);
        Container.RemoveContainerResponse removeResponse = Container.RemoveContainerResponse.newBuilder()
                .setSuccess(false)
                .build();
        when(containerStub.removeContainer(any())).thenReturn(removeResponse);

        assertThrows(dev.whiteo.yadoma.exception.ExecutionConflictException.class, () ->
                containerService.delete(CONTAINER_ID, USER_ID));

        verify(userService).getUserById(USER_ID);
        verify(containerStub).removeContainer(any());
    }

    @Test
    void start_ShouldStartContainer() {
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);
        Container.StartContainerResponse startResponse = Container.StartContainerResponse.newBuilder()
                .setSuccess(true)
                .build();
        when(containerStub.startContainer(any())).thenReturn(startResponse);

        assertDoesNotThrow(() -> containerService.start(CONTAINER_ID, USER_ID));

        verify(userService).isContainerIdContains(CONTAINER_ID, USER_ID);
        verify(containerStub).startContainer(any());
    }

    @Test
    void start_ShouldThrowWhenStartFails() {
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);
        Container.StartContainerResponse startResponse = Container.StartContainerResponse.newBuilder()
                .setSuccess(false)
                .build();
        when(containerStub.startContainer(any())).thenReturn(startResponse);

        assertThrows(dev.whiteo.yadoma.exception.ExecutionConflictException.class, () ->
                containerService.start(CONTAINER_ID, USER_ID));

        verify(containerStub).startContainer(any());
    }

    @Test
    void stop_ShouldStopContainer() {
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);
        Container.StopContainerResponse stopResponse = Container.StopContainerResponse.newBuilder()
                .setSuccess(true)
                .build();
        when(containerStub.stopContainer(any())).thenReturn(stopResponse);

        assertDoesNotThrow(() -> containerService.stop(CONTAINER_ID, USER_ID));

        verify(userService).isContainerIdContains(CONTAINER_ID, USER_ID);
        verify(containerStub).stopContainer(any());
    }

    @Test
    void stop_ShouldThrowWhenStopFails() {
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);
        Container.StopContainerResponse stopResponse = Container.StopContainerResponse.newBuilder()
                .setSuccess(false)
                .build();
        when(containerStub.stopContainer(any())).thenReturn(stopResponse);

        assertThrows(dev.whiteo.yadoma.exception.ExecutionConflictException.class, () ->
                containerService.stop(CONTAINER_ID, USER_ID));

        verify(containerStub).stopContainer(any());
    }

    @Test
    void restart_ShouldRestartContainer() {
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);
        Container.RestartContainerResponse restartResponse = Container.RestartContainerResponse.newBuilder()
                .setSuccess(true)
                .build();
        when(containerStub.restartContainer(any())).thenReturn(restartResponse);

        assertDoesNotThrow(() -> containerService.restart(CONTAINER_ID, USER_ID));

        verify(userService).isContainerIdContains(CONTAINER_ID, USER_ID);
        verify(containerStub).restartContainer(any());
    }

    @Test
    void restart_ShouldThrowWhenRestartFails() {
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);
        Container.RestartContainerResponse restartResponse = Container.RestartContainerResponse.newBuilder()
                .setSuccess(false)
                .build();
        when(containerStub.restartContainer(any())).thenReturn(restartResponse);

        assertThrows(dev.whiteo.yadoma.exception.ExecutionConflictException.class, () ->
                containerService.restart(CONTAINER_ID, USER_ID));

        verify(containerStub).restartContainer(any());
    }

    @Test
    void getById_ShouldHandleGrpcError() {
        when(userService.isContainerIdContains(CONTAINER_ID, USER_ID)).thenReturn(true);
        when(containerStub.getContainerDetails(any())).thenThrow(new RuntimeException("gRPC error"));

        assertThrows(RuntimeException.class, () -> containerService.getById(CONTAINER_ID, USER_ID));

        verify(containerStub).getContainerDetails(any());
    }

    @Test
    void create_ShouldPullImage_whenImageNotExists() {
        ContainerCreateRequest request = new ContainerCreateRequest("test-container", "nonexistent:latest", Collections.emptyList());
        when(userService.getUserById(USER_ID)).thenReturn(testUser);

        image.v1.Image.GetImagesResponse emptyImagesResponse = image.v1.Image.GetImagesResponse.newBuilder().build();
        when(imageStub.getImages(any())).thenReturn(emptyImagesResponse);

        assertThrows(RuntimeException.class, () -> containerService.create(request, USER_ID));

        verify(imageStub).getImages(any());
    }
}
