package dev.whiteo.yadoma.dto.system;

import org.junit.jupiter.api.Test;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class DiskUsageImageResponseTest {

    @Test
    void diskUsageImageResponse_ShouldCreateWithAllFields() {
        String id = "image123";
        List<String> repoTags = List.of("nginx:latest", "nginx:1.21");
        Long size = 100000000L;
        Long containers = 2L;

        DiskUsageImageResponse response = new DiskUsageImageResponse(id, repoTags, size, containers);

        assertEquals(id, response.id());
        assertEquals(repoTags, response.repoTags());
        assertEquals(size, response.size());
        assertEquals(containers, response.containers());
    }

    @Test
    void diskUsageImageResponse_ShouldAllowNullValues() {
        DiskUsageImageResponse response = new DiskUsageImageResponse(null, null, null, null);

        assertNull(response.id());
        assertNull(response.repoTags());
        assertNull(response.size());
        assertNull(response.containers());
    }

    @Test
    void diskUsageImageResponse_ShouldSupportEquality() {
        List<String> tags = List.of("nginx:latest");
        DiskUsageImageResponse response1 = new DiskUsageImageResponse("id1", tags, 100000L, 2L);
        DiskUsageImageResponse response2 = new DiskUsageImageResponse("id1", tags, 100000L, 2L);
        DiskUsageImageResponse response3 = new DiskUsageImageResponse("id2", tags, 100000L, 2L);

        assertEquals(response1, response2);
        assertNotEquals(response1, response3);
        assertEquals(response1.hashCode(), response2.hashCode());
    }

    @Test
    void diskUsageImageResponse_ShouldHaveProperToString() {
        DiskUsageImageResponse response = new DiskUsageImageResponse(
                "id1",
                List.of("nginx:latest"),
                100000L,
                2L
        );

        String toString = response.toString();

        assertNotNull(toString);
        assertTrue(toString.contains("id1"));
        assertTrue(toString.contains("nginx:latest"));
        assertTrue(toString.contains("100000"));
    }
}
