package dev.whiteo.yadoma.dto.system;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class DiskUsageContainerResponseTest {

    @Test
    void diskUsageContainerResponse_ShouldCreateWithAllFields() {
        String id = "container123";
        String image = "nginx:latest";
        String state = "running";
        String status = "Up 5 minutes";
        Long sizeRw = 5000000L;

        DiskUsageContainerResponse response = new DiskUsageContainerResponse(id, image, state, status, sizeRw);

        assertEquals(id, response.id());
        assertEquals(image, response.image());
        assertEquals(state, response.state());
        assertEquals(status, response.status());
        assertEquals(sizeRw, response.sizeRw());
    }

    @Test
    void diskUsageContainerResponse_ShouldAllowNullValues() {
        DiskUsageContainerResponse response = new DiskUsageContainerResponse(null, null, null, null, null);

        assertNull(response.id());
        assertNull(response.image());
        assertNull(response.state());
        assertNull(response.status());
        assertNull(response.sizeRw());
    }

    @Test
    void diskUsageContainerResponse_ShouldSupportEquality() {
        DiskUsageContainerResponse response1 = new DiskUsageContainerResponse(
                "id1", "nginx:latest", "running", "Up 5m", 5000L
        );
        DiskUsageContainerResponse response2 = new DiskUsageContainerResponse(
                "id1", "nginx:latest", "running", "Up 5m", 5000L
        );
        DiskUsageContainerResponse response3 = new DiskUsageContainerResponse(
                "id2", "nginx:latest", "running", "Up 5m", 5000L
        );

        assertEquals(response1, response2);
        assertNotEquals(response1, response3);
        assertEquals(response1.hashCode(), response2.hashCode());
    }

    @Test
    void diskUsageContainerResponse_ShouldHaveProperToString() {
        DiskUsageContainerResponse response = new DiskUsageContainerResponse(
                "id1",
                "nginx:latest",
                "running",
                "Up 5 minutes",
                5000L
        );

        String toString = response.toString();

        assertNotNull(toString);
        assertTrue(toString.contains("id1"));
        assertTrue(toString.contains("nginx:latest"));
        assertTrue(toString.contains("running"));
        assertTrue(toString.contains("Up 5 minutes"));
    }
}
