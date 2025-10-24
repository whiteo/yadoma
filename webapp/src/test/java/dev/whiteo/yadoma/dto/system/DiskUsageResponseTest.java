package dev.whiteo.yadoma.dto.system;

import org.junit.jupiter.api.Test;

import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class DiskUsageResponseTest {

    @Test
    void diskUsageResponse_ShouldCreateWithAllFields() {
        Long layersSize = 1000000L;
        List<DiskUsageImageResponse> images = List.of(
                new DiskUsageImageResponse("img1", List.of("nginx:latest"), 100000L, 2L)
        );
        List<DiskUsageContainerResponse> containers = List.of(
                new DiskUsageContainerResponse("cont1", "nginx:latest", "running", "Up 5m", 5000L)
        );
        List<DiskUsageVolumeResponse> volumes = List.of(
                new DiskUsageVolumeResponse("vol1", "/data", 50000L)
        );

        DiskUsageResponse response = new DiskUsageResponse(layersSize, images, containers, volumes);

        assertEquals(layersSize, response.layersSize());
        assertEquals(images, response.images());
        assertEquals(containers, response.containers());
        assertEquals(volumes, response.volumes());
    }

    @Test
    void diskUsageResponse_ShouldAllowNullValues() {
        DiskUsageResponse response = new DiskUsageResponse(null, null, null, null);

        assertNull(response.layersSize());
        assertNull(response.images());
        assertNull(response.containers());
        assertNull(response.volumes());
    }

    @Test
    void diskUsageResponse_ShouldAllowEmptyLists() {
        DiskUsageResponse response = new DiskUsageResponse(
                0L,
                Collections.emptyList(),
                Collections.emptyList(),
                Collections.emptyList()
        );

        assertEquals(0L, response.layersSize());
        assertTrue(response.images().isEmpty());
        assertTrue(response.containers().isEmpty());
        assertTrue(response.volumes().isEmpty());
    }

    @Test
    void diskUsageResponse_ShouldSupportEquality() {
        List<DiskUsageImageResponse> images = List.of(
                new DiskUsageImageResponse("img1", List.of("nginx:latest"), 100000L, 2L)
        );
        DiskUsageResponse response1 = new DiskUsageResponse(
                1000000L, images, Collections.emptyList(), Collections.emptyList()
        );
        DiskUsageResponse response2 = new DiskUsageResponse(
                1000000L, images, Collections.emptyList(), Collections.emptyList()
        );
        DiskUsageResponse response3 = new DiskUsageResponse(
                2000000L, images, Collections.emptyList(), Collections.emptyList()
        );

        assertEquals(response1, response2);
        assertNotEquals(response1, response3);
        assertEquals(response1.hashCode(), response2.hashCode());
    }

    @Test
    void diskUsageResponse_ShouldHaveProperToString() {
        DiskUsageResponse response = new DiskUsageResponse(
                1000000L,
                Collections.emptyList(),
                Collections.emptyList(),
                Collections.emptyList()
        );

        String toString = response.toString();

        assertNotNull(toString);
        assertTrue(toString.contains("1000000"));
    }
}
