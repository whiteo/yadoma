package dev.whiteo.yadoma.dto.system;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class DiskUsageVolumeResponseTest {

    @Test
    void diskUsageVolumeResponse_ShouldCreateWithAllFields() {
        String name = "volume123";
        String mountpoint = "/var/lib/docker/volumes/volume123";
        Long size = 50000000L;

        DiskUsageVolumeResponse response = new DiskUsageVolumeResponse(name, mountpoint, size);

        assertEquals(name, response.name());
        assertEquals(mountpoint, response.mountpoint());
        assertEquals(size, response.size());
    }

    @Test
    void diskUsageVolumeResponse_ShouldAllowNullValues() {
        DiskUsageVolumeResponse response = new DiskUsageVolumeResponse(null, null, null);

        assertNull(response.name());
        assertNull(response.mountpoint());
        assertNull(response.size());
    }

    @Test
    void diskUsageVolumeResponse_ShouldSupportEquality() {
        DiskUsageVolumeResponse response1 = new DiskUsageVolumeResponse(
                "vol1", "/data/vol1", 50000L
        );
        DiskUsageVolumeResponse response2 = new DiskUsageVolumeResponse(
                "vol1", "/data/vol1", 50000L
        );
        DiskUsageVolumeResponse response3 = new DiskUsageVolumeResponse(
                "vol2", "/data/vol1", 50000L
        );

        assertEquals(response1, response2);
        assertNotEquals(response1, response3);
        assertEquals(response1.hashCode(), response2.hashCode());
    }

    @Test
    void diskUsageVolumeResponse_ShouldHaveProperToString() {
        DiskUsageVolumeResponse response = new DiskUsageVolumeResponse(
                "vol1",
                "/var/lib/docker/volumes/vol1",
                50000L
        );

        String toString = response.toString();

        assertNotNull(toString);
        assertTrue(toString.contains("vol1"));
        assertTrue(toString.contains("/var/lib/docker/volumes/vol1"));
        assertTrue(toString.contains("50000"));
    }
}
