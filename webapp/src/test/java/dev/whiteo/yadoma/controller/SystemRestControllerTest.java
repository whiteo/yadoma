package dev.whiteo.yadoma.controller;

import dev.whiteo.yadoma.dto.system.DiskUsageResponse;
import dev.whiteo.yadoma.dto.system.SystemInfoResponse;
import dev.whiteo.yadoma.service.SystemService;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class SystemRestControllerTest {

    @Mock
    private SystemService systemService;

    @InjectMocks
    private SystemRestController systemRestController;

    @Test
    void getSystemInfo_ShouldReturnSystemInfo() {
        SystemInfoResponse systemInfo = new SystemInfoResponse(
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

        when(systemService.getSystemInfo()).thenReturn(systemInfo);

        ResponseEntity<SystemInfoResponse> response = systemRestController.getSystemInfo();

        assertEquals(HttpStatus.OK, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(systemInfo, response.getBody());
        assertEquals("test-id", response.getBody().id());
        assertEquals("test-name", response.getBody().name());
        assertEquals("27.0.0", response.getBody().serverVersion());
        assertEquals(8, response.getBody().nCpu());
        assertEquals(16000000000L, response.getBody().memTotal());
        verify(systemService).getSystemInfo();
    }

    @Test
    void getDiskUsage_ShouldReturnDiskUsage() {
        DiskUsageResponse diskUsage = new DiskUsageResponse(
                1000000L,
                Collections.emptyList(),
                Collections.emptyList(),
                Collections.emptyList()
        );

        when(systemService.getDiskUsage()).thenReturn(diskUsage);

        ResponseEntity<DiskUsageResponse> response = systemRestController.getDiskUsage();

        assertEquals(HttpStatus.OK, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(diskUsage, response.getBody());
        assertEquals(1000000L, response.getBody().layersSize());
        verify(systemService).getDiskUsage();
    }
}
