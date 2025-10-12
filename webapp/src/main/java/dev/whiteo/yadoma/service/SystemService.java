package dev.whiteo.yadoma.service;

import dev.whiteo.yadoma.dto.system.DiskUsageResponse;
import dev.whiteo.yadoma.dto.system.SystemInfoResponse;
import dev.whiteo.yadoma.mapper.SystemMapper;
import io.grpc.StatusRuntimeException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import system.v1.System;
import system.v1.SystemServiceGrpc;

import java.util.concurrent.TimeUnit;

/**
 * Service for managing system information and disk usage.
 * Provides methods to retrieve Docker system info and disk usage via gRPC.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class SystemService {

    private final SystemServiceGrpc.SystemServiceBlockingStub systemStub;
    private final SystemMapper mapper;

    /**
     * Retrieves system information from Docker daemon.
     *
     * @return SystemInfoResponse containing system details
     * @throws RuntimeException if gRPC agent is unavailable or connection fails
     */
    public SystemInfoResponse getSystemInfo() {
        try {
            System.GetSystemInfoRequest request = System.GetSystemInfoRequest.newBuilder().build();
            System.GetSystemInfoResponse response = systemStub
                    .withDeadlineAfter(30, TimeUnit.SECONDS)
                    .getSystemInfo(request);
            return mapper.toSystemInfoResponse(response);
        } catch (StatusRuntimeException e) {
            log.error("gRPC error while getting system info: {} - {}", e.getStatus(), e.getMessage());
            throw new RuntimeException("Agent service unavailable. Please ensure the Yadoma agent is running and accessible.", e);
        } catch (Exception e) {
            log.error("Unexpected error while getting system info", e);
            throw new RuntimeException("Failed to get system info: " + e.getMessage(), e);
        }
    }

    /**
     * Retrieves disk usage information for images, containers, and volumes.
     *
     * @return DiskUsageResponse containing disk usage details
     * @throws RuntimeException if gRPC agent is unavailable or connection fails
     */
    public DiskUsageResponse getDiskUsage() {
        try {
            System.GetDiskUsageRequest request = System.GetDiskUsageRequest.newBuilder().build();
            System.GetDiskUsageResponse response = systemStub
                    .withDeadlineAfter(10, TimeUnit.SECONDS)
                    .getDiskUsage(request);
            return mapper.toDiskUsageResponse(response);
        } catch (StatusRuntimeException e) {
            log.error("gRPC error while getting disk usage: {} - {}", e.getStatus(), e.getMessage());
            throw new RuntimeException("Agent service unavailable. Please ensure the Yadoma agent is running and accessible.", e);
        } catch (Exception e) {
            log.error("Unexpected error while getting disk usage", e);
            throw new RuntimeException("Failed to get disk usage: " + e.getMessage(), e);
        }
    }
}
