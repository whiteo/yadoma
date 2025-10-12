package dev.whiteo.yadoma.service;

import container.v1.Container;
import container.v1.ContainerServiceGrpc;
import dev.whiteo.yadoma.domain.User;
import dev.whiteo.yadoma.dto.container.ContainerCreateRequest;
import dev.whiteo.yadoma.dto.container.ContainerResponse;
import dev.whiteo.yadoma.exception.ExecutionConflictException;
import dev.whiteo.yadoma.mapper.ContainerMapper;
import image.v1.Image;
import image.v1.ImageServiceGrpc;
import io.grpc.StatusRuntimeException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.stereotype.Service;

import java.util.Iterator;
import java.util.List;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;

/**
 * Service for managing containers.
 * Provides methods to find, create, delete, start, stop and restart containers.
 * Validates user access to containers.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class ContainerService {

    private final ContainerServiceGrpc.ContainerServiceBlockingStub containerStub;
    private final ImageServiceGrpc.ImageServiceBlockingStub imageStub;
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

        log.info("Finding containers for user {}. User has {} container IDs: {}",
                userId, user.getContainerIds() == null ? 0 : user.getContainerIds().size(), user.getContainerIds());

        List<ContainerResponse> containerResponses = null;
        try {
            Container.GetContainersRequest request = Container.GetContainersRequest
                    .newBuilder()
                    .setAll(true)
                    .build();
            Container.GetContainersResponse response = containerStub.getContainers(request);
            log.info("Got {} containers from gRPC", response.getContainersList().size());

            containerResponses = response.getContainersList()
                    .stream()
                    .filter(container -> user.getContainerIds() != null && user.getContainerIds().contains(container.getId()))
                    .map(mapper::toResponseDTO)
                    .collect(Collectors.toList());

            log.info("After filtering, {} containers remain", containerResponses.size());
        } catch (Exception e) {
            log.error("Error fetching containers", e);
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
     * Automatically pulls the Docker image if it doesn't exist locally.
     *
     * @param request the container creation request
     * @param userId  the ID of the user
     * @throws ExecutionConflictException if the container creation fails
     */
    public void create(ContainerCreateRequest request, String userId) {
        User user = userService.getUserById(userId);

        try {
            ensureImageExists(request.image());

            Container.CreateContainerRequest grpcRequest = Container.CreateContainerRequest
                    .newBuilder()
                    .setImage(request.image())
                    .setName(request.name())
                    .addAllEnv(request.envVars())
                    .build();
            Container.CreateContainerResponse response = containerStub.createContainer(grpcRequest);
            userService.addContainerToUser(response.getId(), user);

            log.info("Container created successfully: {} with image: {}", response.getId(), request.image());
        } catch (StatusRuntimeException e) {
            log.error("gRPC error creating container: {}", e.getStatus());
            throw new ExecutionConflictException("Failed to create container: " + e.getStatus().getDescription());
        } catch (Exception e) {
            log.error("Unexpected error creating container", e);
            throw new ExecutionConflictException("Failed to create container: " + e.getMessage());
        }
    }

    /**
     * Ensures a Docker image exists locally.
     * If the image doesn't exist, it will be pulled from the registry.
     *
     * @param imageName the name of the image (e.g., "nginx:latest", "redis:7.0")
     */
    private void ensureImageExists(String imageName) {
        try {
            if (imageExists(imageName)) {
                log.info("Image already exists locally: {}", imageName);
                return;
            }

            log.info("Image not found locally, pulling: {}", imageName);
            pullImage(imageName);
            log.info("Image pulled successfully: {}", imageName);

        } catch (Exception e) {
            log.error("Failed to ensure image exists: {}", imageName, e);
            throw new ExecutionConflictException("Failed to pull image: " + imageName);
        }
    }

    /**
     * Checks if a Docker image exists locally.
     *
     * @param imageName the name of the image
     * @return true if the image exists, false otherwise
     */
    private boolean imageExists(String imageName) {
        try {
            Image.GetImagesRequest request = Image.GetImagesRequest.newBuilder()
                    .setAll(true)
                    .build();
            Image.GetImagesResponse response = imageStub.getImages(request);

            return response.getImagesList().stream()
                    .flatMap(image -> image.getRepoTagsList().stream())
                    .anyMatch(tag -> matchesImageName(tag, imageName));
        } catch (Exception e) {
            log.warn("Error checking if image exists: {}", imageName, e);
            return false;
        }
    }

    /**
     * Checks if a repo tag matches the requested image name.
     * Handles cases like "nginx" matching "nginx:latest" and vice versa.
     *
     * @param repoTag   the repository tag (e.g., "nginx:latest")
     * @param imageName the requested image name (e.g., "nginx" or "nginx:latest")
     * @return true if they match
     */
    private boolean matchesImageName(String repoTag, String imageName) {
        if (repoTag.equals(imageName)) {
            return true;
        }

        String imageWithTag = imageName.contains(":") ? imageName : imageName + ":latest";
        return repoTag.equals(imageWithTag);
    }

    /**
     * Pulls a Docker image from the registry.
     *
     * @param imageName the name of the image to pull
     */
    private void pullImage(String imageName) {
        try {
            String imageLink = imageName.contains(":") ? imageName : imageName + ":latest";

            Image.PullImageRequest request = Image.PullImageRequest.newBuilder()
                    .setLink(imageLink)
                    .build();

            log.info("Starting image pull for: {})", imageLink);

            Iterator<Image.PullImageResponse> responseIterator = imageStub
                    .withDeadlineAfter(10, TimeUnit.SECONDS)
                    .pullImage(request);

            int chunkCount = 0;
            long totalBytes = 0;

            try {
                while (responseIterator.hasNext()) {
                    Image.PullImageResponse chunk = responseIterator.next();
                    chunkCount++;
                    totalBytes += chunk.getChunk().size();

                    if (chunkCount % 10 == 0) {
                        log.info("Pulling image {}: {} chunks received, {} bytes total",
                                imageLink, chunkCount, totalBytes);
                    }
                }
            } catch (Exception e) {
                log.error("Error during streaming image pull: {}", e.getMessage());
                throw e;
            }

            log.info("Image pull completed: {} ({} chunks, {} total bytes)",
                    imageLink, chunkCount, totalBytes);

        } catch (StatusRuntimeException e) {
            log.error("gRPC error pulling image: {} - {}", e.getStatus().getCode(), e.getStatus().getDescription());

            String errorMessage = switch (e.getStatus().getCode()) {
                case CANCELLED -> "Image pull was cancelled. The operation may have timed out or been interrupted.";
                case DEADLINE_EXCEEDED -> "Image pull timed out. The image may be too large or network is too slow.";
                case NOT_FOUND -> "Image not found: " + imageName + ". Please verify the image name and tag.";
                case UNAVAILABLE -> "Docker registry is unavailable. Please check your internet connection.";
                case PERMISSION_DENIED -> "Permission denied. You may need to authenticate with the registry.";
                default -> "Failed to pull image: " + e.getStatus().getDescription();
            };

            throw new ExecutionConflictException(errorMessage);
        } catch (Exception e) {
            log.error("Unexpected error pulling image: {}", imageName, e);
            throw new ExecutionConflictException("Unexpected error pulling image: " + e.getMessage());
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
     * Gets container logs as a stream.
     *
     * @param containerId the ID of the container
     * @param userId      the ID of the user
     * @param follow      whether to follow the log stream
     * @return iterator of log chunks
     * @throws BadCredentialsException if access denied
     */
    public Iterator<Container.GetContainerLogsResponse> getLogs(String containerId, String userId, boolean follow) {
        validateUserAccess(containerId, userId);

        Container.GetContainerLogsRequest request = Container.GetContainerLogsRequest
                .newBuilder()
                .setId(containerId)
                .setFollow(follow)
                .build();

        return containerStub.getContainerLogs(request);
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