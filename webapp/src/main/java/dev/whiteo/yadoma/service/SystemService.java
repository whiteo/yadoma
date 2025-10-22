package dev.whiteo.yadoma.service;

import dev.whiteo.yadoma.dto.system.DiskUsageResponse;
import dev.whiteo.yadoma.dto.system.SystemInfoResponse;
import dev.whiteo.yadoma.mapper.SystemMapper;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import system.v1.System;
import system.v1.SystemServiceGrpc;

/**
 * Service for managing system information and disk usage.
 * Provides methods to retrieve Docker system info and disk usage via gRPC.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Service
@RequiredArgsConstructor
public class SystemService {

    private final SystemServiceGrpc.SystemServiceBlockingStub systemStub;
    private final SystemMapper mapper;

    /**
     * Retrieves system information from Docker daemon.
     *
     * @return SystemInfoResponse containing system details
     */
    public SystemInfoResponse getSystemInfo() {
        try {
            System.GetSystemInfoRequest request = System.GetSystemInfoRequest.newBuilder().build();
            System.GetSystemInfoResponse response = systemStub.getSystemInfo(request);
            return mapper.toSystemInfoResponse(response);
        } catch (Exception e) {
            throw new RuntimeException("Failed to get system info: " + e.getMessage(), e);
        }
    }

    /**
     * Retrieves disk usage information for images, containers, and volumes.
     *
     * @return DiskUsageResponse containing disk usage details
     */
    public DiskUsageResponse getDiskUsage() {
        try {
            System.GetDiskUsageRequest request = System.GetDiskUsageRequest.newBuilder().build();
            System.GetDiskUsageResponse response = systemStub.getDiskUsage(request);
            return mapper.toDiskUsageResponse(response);
        } catch (Exception e) {
            throw new RuntimeException("Failed to get disk usage: " + e.getMessage(), e);
        }
    }
}
