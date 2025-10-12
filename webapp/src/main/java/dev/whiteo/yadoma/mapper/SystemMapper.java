package dev.whiteo.yadoma.mapper;

import dev.whiteo.yadoma.dto.system.*;
import org.mapstruct.Mapper;
import system.v1.System;

import java.util.List;
import java.util.stream.Collectors;

/**
 * Mapper for converting System gRPC responses to DTOs.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Mapper(componentModel = "spring")
public interface SystemMapper {

    /**
     * Converts GetSystemInfoResponse to SystemInfoResponse DTO.
     *
     * @param response gRPC response
     * @return SystemInfoResponse DTO
     */
    default SystemInfoResponse toSystemInfoResponse(System.GetSystemInfoResponse response) {
        return new SystemInfoResponse(
                response.getId(),
                response.getName(),
                response.getServerVersion(),
                response.getKernelVersion(),
                response.getOperatingSystem(),
                response.getArchitecture(),
                response.getNCpu(),
                response.getMemTotal(),
                response.getContainers(),
                response.getContainersRunning(),
                response.getContainersPaused(),
                response.getContainersStopped(),
                response.getImages(),
                response.getDriver(),
                response.getLabelsList()
        );
    }

    /**
     * Converts GetDiskUsageResponse to DiskUsageResponse DTO.
     *
     * @param response gRPC response
     * @return DiskUsageResponse DTO
     */
    default DiskUsageResponse toDiskUsageResponse(System.GetDiskUsageResponse response) {
        List<DiskUsageImageResponse> images = response.getImagesList().stream()
                .map(this::toDiskUsageImageResponse)
                .collect(Collectors.toList());

        List<DiskUsageContainerResponse> containers = response.getContainersList().stream()
                .map(this::toDiskUsageContainerResponse)
                .collect(Collectors.toList());

        List<DiskUsageVolumeResponse> volumes = response.getVolumesList().stream()
                .map(this::toDiskUsageVolumeResponse)
                .collect(Collectors.toList());

        return new DiskUsageResponse(
                response.getLayersSize(),
                images,
                containers,
                volumes
        );
    }

    /**
     * Converts DiskUsageImage to DiskUsageImageResponse DTO.
     *
     * @param image gRPC disk usage image
     * @return DiskUsageImageResponse DTO
     */
    default DiskUsageImageResponse toDiskUsageImageResponse(System.DiskUsageImage image) {
        return new DiskUsageImageResponse(
                image.getId(),
                image.getRepoTagsList(),
                image.getSize(),
                image.getContainers()
        );
    }

    /**
     * Converts DiskUsageContainer to DiskUsageContainerResponse DTO.
     *
     * @param container gRPC disk usage container
     * @return DiskUsageContainerResponse DTO
     */
    default DiskUsageContainerResponse toDiskUsageContainerResponse(System.DiskUsageContainer container) {
        return new DiskUsageContainerResponse(
                container.getId(),
                container.getImage(),
                container.getState(),
                container.getStatus(),
                container.getSizeRw()
        );
    }

    /**
     * Converts DiskUsageVolume to DiskUsageVolumeResponse DTO.
     *
     * @param volume gRPC disk usage volume
     * @return DiskUsageVolumeResponse DTO
     */
    default DiskUsageVolumeResponse toDiskUsageVolumeResponse(System.DiskUsageVolume volume) {
        return new DiskUsageVolumeResponse(
                volume.getName(),
                volume.getMountpoint(),
                volume.getSize()
        );
    }
}
