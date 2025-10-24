package dev.whiteo.yadoma.service;

import dev.whiteo.yadoma.dto.system.DiskUsageResponse;
import dev.whiteo.yadoma.dto.system.SystemInfoResponse;
import dev.whiteo.yadoma.mapper.SystemMapper;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import system.v1.System;
import system.v1.SystemServiceGrpc;

import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class SystemServiceTest {

    @Mock
    private SystemServiceGrpc.SystemServiceBlockingStub systemStub;

    @Mock
    private SystemMapper mapper;

    @InjectMocks
    private SystemService systemService;

    private System.GetSystemInfoResponse grpcSystemInfoResponse;
    private SystemInfoResponse systemInfoResponse;
    private System.GetDiskUsageResponse grpcDiskUsageResponse;
    private DiskUsageResponse diskUsageResponse;

    @BeforeEach
    void setUp() {
        grpcSystemInfoResponse = System.GetSystemInfoResponse.newBuilder()
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

        systemInfoResponse = new SystemInfoResponse(
                "test-id",
                "test-name",
                "27.0.0",
                "6.8.0-50-generic",
                "Ubuntu 24.04 LTS",
                "x86_64",
                8,
                16000000000L,
                10,
                3,
                0,
                7,
                5,
                "overlay2",
                List.of("label1", "label2")
        );

        grpcDiskUsageResponse = System.GetDiskUsageResponse.newBuilder()
                .setLayersSize(1000000L)
                .build();

        diskUsageResponse = new DiskUsageResponse(
                1000000L,
                Collections.emptyList(),
                Collections.emptyList(),
                Collections.emptyList()
        );
    }

    @Test
    void getSystemInfo_ShouldReturnSystemInfo() {
        when(systemStub.withDeadlineAfter(anyLong(), any())).thenReturn(systemStub);
        when(systemStub.getSystemInfo(any(System.GetSystemInfoRequest.class))).thenReturn(grpcSystemInfoResponse);
        when(mapper.toSystemInfoResponse(grpcSystemInfoResponse)).thenReturn(systemInfoResponse);

        SystemInfoResponse result = systemService.getSystemInfo();

        assertNotNull(result);
        assertEquals(systemInfoResponse, result);
        assertEquals("test-id", result.id());
        assertEquals("test-name", result.name());
        assertEquals("27.0.0", result.serverVersion());
        verify(systemStub).withDeadlineAfter(30, java.util.concurrent.TimeUnit.SECONDS);
        verify(systemStub).getSystemInfo(any(System.GetSystemInfoRequest.class));
        verify(mapper).toSystemInfoResponse(grpcSystemInfoResponse);
    }

    @Test
    void getSystemInfo_ShouldThrowRuntimeExceptionOnGrpcError() {
        when(systemStub.withDeadlineAfter(anyLong(), any())).thenReturn(systemStub);
        when(systemStub.getSystemInfo(any(System.GetSystemInfoRequest.class)))
                .thenThrow(new StatusRuntimeException(Status.UNAVAILABLE));

        RuntimeException exception = assertThrows(RuntimeException.class, () ->
                systemService.getSystemInfo());

        assertTrue(exception.getMessage().contains("Agent service unavailable"));
        verify(systemStub).getSystemInfo(any(System.GetSystemInfoRequest.class));
        verifyNoInteractions(mapper);
    }

    @Test
    void getSystemInfo_ShouldThrowRuntimeExceptionOnUnexpectedError() {
        when(systemStub.withDeadlineAfter(anyLong(), any())).thenReturn(systemStub);
        when(systemStub.getSystemInfo(any(System.GetSystemInfoRequest.class)))
                .thenThrow(new RuntimeException("Unexpected error"));

        RuntimeException exception = assertThrows(RuntimeException.class, () ->
                systemService.getSystemInfo());

        assertTrue(exception.getMessage().contains("Failed to get system info"));
        verify(systemStub).getSystemInfo(any(System.GetSystemInfoRequest.class));
        verifyNoInteractions(mapper);
    }

    @Test
    void getDiskUsage_ShouldReturnDiskUsage() {
        when(systemStub.withDeadlineAfter(anyLong(), any())).thenReturn(systemStub);
        when(systemStub.getDiskUsage(any(System.GetDiskUsageRequest.class))).thenReturn(grpcDiskUsageResponse);
        when(mapper.toDiskUsageResponse(grpcDiskUsageResponse)).thenReturn(diskUsageResponse);

        DiskUsageResponse result = systemService.getDiskUsage();

        assertNotNull(result);
        assertEquals(diskUsageResponse, result);
        assertEquals(1000000L, result.layersSize());
        verify(systemStub).withDeadlineAfter(10, java.util.concurrent.TimeUnit.SECONDS);
        verify(systemStub).getDiskUsage(any(System.GetDiskUsageRequest.class));
        verify(mapper).toDiskUsageResponse(grpcDiskUsageResponse);
    }

    @Test
    void getDiskUsage_ShouldThrowRuntimeExceptionOnGrpcError() {
        when(systemStub.withDeadlineAfter(anyLong(), any())).thenReturn(systemStub);
        when(systemStub.getDiskUsage(any(System.GetDiskUsageRequest.class)))
                .thenThrow(new StatusRuntimeException(Status.UNAVAILABLE));

        RuntimeException exception = assertThrows(RuntimeException.class, () ->
                systemService.getDiskUsage());

        assertTrue(exception.getMessage().contains("Agent service unavailable"));
        verify(systemStub).getDiskUsage(any(System.GetDiskUsageRequest.class));
        verifyNoInteractions(mapper);
    }

    @Test
    void getDiskUsage_ShouldThrowRuntimeExceptionOnUnexpectedError() {
        when(systemStub.withDeadlineAfter(anyLong(), any())).thenReturn(systemStub);
        when(systemStub.getDiskUsage(any(System.GetDiskUsageRequest.class)))
                .thenThrow(new RuntimeException("Unexpected error"));

        RuntimeException exception = assertThrows(RuntimeException.class, () ->
                systemService.getDiskUsage());

        assertTrue(exception.getMessage().contains("Failed to get disk usage"));
        verify(systemStub).getDiskUsage(any(System.GetDiskUsageRequest.class));
        verifyNoInteractions(mapper);
    }
}
