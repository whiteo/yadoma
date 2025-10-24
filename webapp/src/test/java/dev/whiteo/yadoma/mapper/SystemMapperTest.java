package dev.whiteo.yadoma.mapper;

import dev.whiteo.yadoma.dto.system.*;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mapstruct.factory.Mappers;
import system.v1.System;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class SystemMapperTest {

    private SystemMapper systemMapper;

    @BeforeEach
    void setUp() {
        systemMapper = Mappers.getMapper(SystemMapper.class);
    }

    @Test
    void toSystemInfoResponse_ShouldMapCompleteResponse() {
        System.GetSystemInfoResponse grpcResponse = System.GetSystemInfoResponse.newBuilder()
                .setId("test-id")
                .setName("test-name")
                .setServerVersion("27.0.0")
                .setKernelVersion("6.8.0-50-generic")
                .setOperatingSystem("Ubuntu 24.04 LTS")
                .setArchitecture("x86_64")
                .setNCpu(8)
                .setMemTotal(16000000000L)
                .setContainers(10)
                .setContainersRunning(3)
                .setContainersPaused(0)
                .setContainersStopped(7)
                .setImages(5)
                .setDriver("overlay2")
                .addLabels("label1")
                .addLabels("label2")
                .build();

        SystemInfoResponse result = systemMapper.toSystemInfoResponse(grpcResponse);

        assertNotNull(result);
        assertEquals("test-id", result.id());
        assertEquals("test-name", result.name());
        assertEquals("27.0.0", result.serverVersion());
        assertEquals("6.8.0-50-generic", result.kernelVersion());
        assertEquals("Ubuntu 24.04 LTS", result.operatingSystem());
        assertEquals("x86_64", result.architecture());
        assertEquals(8, result.nCpu());
        assertEquals(16000000000L, result.memTotal());
        assertEquals(10, result.containers());
        assertEquals(3, result.containersRunning());
        assertEquals(0, result.containersPaused());
        assertEquals(7, result.containersStopped());
        assertEquals(5, result.images());
        assertEquals("overlay2", result.driver());
        assertEquals(List.of("label1", "label2"), result.labels());
    }

    @Test
    void toSystemInfoResponse_ShouldHandleEmptyLabels() {
        System.GetSystemInfoResponse grpcResponse = System.GetSystemInfoResponse.newBuilder()
                .setId("test-id")
                .setName("test-name")
                .setServerVersion("27.0.0")
                .setKernelVersion("6.8.0")
                .setOperatingSystem("Ubuntu")
                .setArchitecture("x86_64")
                .setNCpu(4)
                .setMemTotal(8000000000L)
                .setContainers(0)
                .setContainersRunning(0)
                .setContainersPaused(0)
                .setContainersStopped(0)
                .setImages(0)
                .setDriver("overlay2")
                .build();

        SystemInfoResponse result = systemMapper.toSystemInfoResponse(grpcResponse);

        assertNotNull(result);
        assertTrue(result.labels().isEmpty());
    }

    @Test
    void toDiskUsageResponse_ShouldMapCompleteResponse() {
        System.DiskUsageImage image = System.DiskUsageImage.newBuilder()
                .setId("image123")
                .addRepoTags("nginx:latest")
                .setSize(100000000L)
                .setContainers(2)
                .build();

        System.DiskUsageContainer container = System.DiskUsageContainer.newBuilder()
                .setId("container123")
                .setImage("nginx:latest")
                .setState("running")
                .setStatus("Up 5 minutes")
                .setSizeRw(5000000L)
                .build();

        System.DiskUsageVolume volume = System.DiskUsageVolume.newBuilder()
                .setName("volume123")
                .setMountpoint("/var/lib/docker/volumes/volume123")
                .setSize(50000000L)
                .build();

        System.GetDiskUsageResponse grpcResponse = System.GetDiskUsageResponse.newBuilder()
                .setLayersSize(1000000L)
                .addImages(image)
                .addContainers(container)
                .addVolumes(volume)
                .build();

        DiskUsageResponse result = systemMapper.toDiskUsageResponse(grpcResponse);

        assertNotNull(result);
        assertEquals(1000000L, result.layersSize());
        assertEquals(1, result.images().size());
        assertEquals(1, result.containers().size());
        assertEquals(1, result.volumes().size());
    }

    @Test
    void toDiskUsageResponse_ShouldHandleEmptyLists() {
        System.GetDiskUsageResponse grpcResponse = System.GetDiskUsageResponse.newBuilder()
                .setLayersSize(0L)
                .build();

        DiskUsageResponse result = systemMapper.toDiskUsageResponse(grpcResponse);

        assertNotNull(result);
        assertEquals(0L, result.layersSize());
        assertTrue(result.images().isEmpty());
        assertTrue(result.containers().isEmpty());
        assertTrue(result.volumes().isEmpty());
    }

    @Test
    void toDiskUsageImageResponse_ShouldMapCorrectly() {
        System.DiskUsageImage grpcImage = System.DiskUsageImage.newBuilder()
                .setId("image123")
                .addRepoTags("nginx:latest")
                .addRepoTags("nginx:1.21")
                .setSize(100000000L)
                .setContainers(2)
                .build();

        DiskUsageImageResponse result = systemMapper.toDiskUsageImageResponse(grpcImage);

        assertNotNull(result);
        assertEquals("image123", result.id());
        assertEquals(List.of("nginx:latest", "nginx:1.21"), result.repoTags());
        assertEquals(100000000L, result.size());
        assertEquals(2, result.containers());
    }

    @Test
    void toDiskUsageContainerResponse_ShouldMapCorrectly() {
        System.DiskUsageContainer grpcContainer = System.DiskUsageContainer.newBuilder()
                .setId("container123")
                .setImage("nginx:latest")
                .setState("running")
                .setStatus("Up 5 minutes")
                .setSizeRw(5000000L)
                .build();

        DiskUsageContainerResponse result = systemMapper.toDiskUsageContainerResponse(grpcContainer);

        assertNotNull(result);
        assertEquals("container123", result.id());
        assertEquals("nginx:latest", result.image());
        assertEquals("running", result.state());
        assertEquals("Up 5 minutes", result.status());
        assertEquals(5000000L, result.sizeRw());
    }

    @Test
    void toDiskUsageVolumeResponse_ShouldMapCorrectly() {
        System.DiskUsageVolume grpcVolume = System.DiskUsageVolume.newBuilder()
                .setName("volume123")
                .setMountpoint("/var/lib/docker/volumes/volume123")
                .setSize(50000000L)
                .build();

        DiskUsageVolumeResponse result = systemMapper.toDiskUsageVolumeResponse(grpcVolume);

        assertNotNull(result);
        assertEquals("volume123", result.name());
        assertEquals("/var/lib/docker/volumes/volume123", result.mountpoint());
        assertEquals(50000000L, result.size());
    }
}
