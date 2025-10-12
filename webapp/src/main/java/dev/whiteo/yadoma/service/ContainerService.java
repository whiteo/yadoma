package dev.whiteo.yadoma.service;

import container.v1.Container;
import container.v1.ContainerServiceGrpc;
import dev.whiteo.yadoma.domain.User;
import dev.whiteo.yadoma.dto.container.ContainerCreateRequest;
import dev.whiteo.yadoma.dto.container.ContainerResponse;
import dev.whiteo.yadoma.exception.ExecutionConflictException;
import dev.whiteo.yadoma.mapper.ContainerMapper;
import lombok.RequiredArgsConstructor;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.stream.Collectors;

/**
 * Service for managing containers.
 * Provides methods to find, create, delete, start, stop and restart containers.
 * Validates user access to containers.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Service
@RequiredArgsConstructor
public class ContainerService {

    private final ContainerServiceGrpc.ContainerServiceBlockingStub containerStub;
    private final UserService userService;
    private final ContainerMapper mapper;

    /**
     * Finds all containers accessible by the specified user.
     *
     * @param userId  the ID of the user
     * @param adminId the ID of the admin
     * @return a list of ContainerResponse DTOs representing the accessible containers
     */
    public List<ContainerResponse> findAll(String userId, String adminId) {
        User user = userService.validateUserAccess(userId, adminId);

        List<ContainerResponse> containerResponses = null;
        try {
            Container.GetContainersRequest request = Container.GetContainersRequest
                    .newBuilder()
                    .setAll(true)
                    .build();
            Container.GetContainersResponse response = containerStub.getContainers(request);
            response.getContainersList().removeIf(container ->
                    !user.getContainerIds().contains(container.getId())
            );
            containerResponses = response.getContainersList()
                    .stream()
                    .map(mapper::toResponseDTO)
                    .collect(Collectors.toList());
        } catch (Exception e) {
            System.out.println(e.getMessage());
        }

        return containerResponses;
    }

    /**
     * Retrieves a container by its ID.
     *
     * @param containerId the ID of the container
     * @param userId      the ID of the user
     * @return the ContainerResponse DTO for the container
     * @throws ExecutionConflictException if the container retrieval fails
     */
    public ContainerResponse getById(String containerId, String userId) {
        validateUserAccess(containerId, userId);

        ContainerResponse containerResponse = null;
        try {
            Container.GetContainerDetailsRequest request = Container.GetContainerDetailsRequest
                    .newBuilder()
                    .setId(containerId)
                    .build();
            Container.GetContainerDetailsResponse response = containerStub.getContainerDetails(request);
            containerResponse = mapper.toResponseDTO(response);
        } catch (Exception e) {
            System.out.println(e.getMessage());
        }

        return containerResponse;
    }

    /**
     * Creates a new container.
     *
     * @param request the container creation request
     * @param userId  the ID of the user
     * @throws ExecutionConflictException if the container creation fails
     */
    public void create(ContainerCreateRequest request, String userId) {
        User user = userService.getUserById(userId);

        try {
            Container.CreateContainerRequest grpcRequest = Container.CreateContainerRequest
                    .newBuilder()
                    .setImage(request.image())
                    .setName(request.name())
                    .addAllEnv(request.envVars())
                    .build();
            Container.CreateContainerResponse response = containerStub.createContainer(grpcRequest);
            userService.addContainerToUser(response.getId(), user);
        } catch (Exception e) {
            System.out.println(e.getMessage());
        }
    }

    /**
     * Deletes a container.
     *
     * @param containerId the ID of the container
     * @param userId      the ID of the user
     * @throws ExecutionConflictException if the container deletion fails
     */
    public void delete(String containerId, String userId) {
        User user = userService.getUserById(userId);
        validateUserAccess(containerId, user.getId());

        Container.RemoveContainerRequest request = Container.RemoveContainerRequest
                .newBuilder()
                .setId(containerId)
                .setForce(true)
                .setRemoveVolumes(true)
                .build();
        Container.RemoveContainerResponse response = containerStub.removeContainer(request);
        if (!response.getSuccess()) {
            throw new ExecutionConflictException("Failed to remove container: " + containerId);
        }

        userService.removeContainerFromUser(containerId, user);
    }

    /**
     * Starts a container.
     *
     * @param containerId the ID of the container
     * @param userId      the ID of the user
     * @throws ExecutionConflictException if the container start fails
     */
    public void start(String containerId, String userId) {
        validateUserAccess(containerId, userId);

        Container.StartContainerRequest request = Container.StartContainerRequest
                .newBuilder()
                .setId(containerId)
                .build();
        Container.StartContainerResponse response = containerStub.startContainer(request);
        if (!response.getSuccess()) {
            throw new ExecutionConflictException("Failed to start container: " + containerId);
        }
    }

    /**
     * Stops a container.
     *
     * @param containerId the ID of the container
     * @param userId      the ID of the user
     * @throws ExecutionConflictException if the container stop fails
     */
    public void stop(String containerId, String userId) {
        validateUserAccess(containerId, userId);

        Container.StopContainerRequest request = Container.StopContainerRequest
                .newBuilder()
                .setId(containerId)
                .build();
        Container.StopContainerResponse response = containerStub.stopContainer(request);
        if (!response.getSuccess()) {
            throw new ExecutionConflictException("Failed to stop container: " + containerId);
        }
    }

    /**
     * Restarts a container.
     *
     * @param containerId the ID of the container
     * @param userId      the ID of the user
     * @throws ExecutionConflictException if the container restart fails
     */
    public void restart(String containerId, String userId) {
        validateUserAccess(containerId, userId);

        Container.RestartContainerRequest request = Container.RestartContainerRequest
                .newBuilder()
                .setId(containerId)
                .build();
        Container.RestartContainerResponse response = containerStub.restartContainer(request);
        if (!response.getSuccess()) {
            throw new ExecutionConflictException("Failed to restart container: " + containerId);
        }
    }

    /**
     * Validates user access to a container.
     *
     * @param containerId the ID of the container
     * @param userId      the ID of the user
     * @throws BadCredentialsException if the container ID does not belong to the user
     */
    private void validateUserAccess(String containerId, String userId) {
        boolean contains = userService.isContainerIdContains(containerId, userId);
        if (!contains) {
            throw new BadCredentialsException("Container ID does not belong to the user");
        }
    }
}