package dev.whiteo.yadoma.dto.system;

import org.junit.jupiter.api.Test;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class SystemInfoResponseTest {

    @Test
    void systemInfoResponse_ShouldCreateWithAllFields() {
        String id = "test-id";
        String name = "test-name";
        String serverVersion = "27.0.0";
        String kernelVersion = "6.8.0-50-generic";
        String operatingSystem = "Ubuntu 24.04 LTS";
        String architecture = "x86_64";
        Integer nCpu = 8;
        Long memTotal = 16000000000L;
        Integer containers = 10;
        Integer containersRunning = 3;
        Integer containersPaused = 0;
        Integer containersStopped = 7;
        Integer images = 5;
        String driver = "overlay2";
        List<String> labels = List.of("label1", "label2");

        SystemInfoResponse response = new SystemInfoResponse(
                id, name, serverVersion, kernelVersion, operatingSystem,
                architecture, nCpu, memTotal, containers, containersRunning,
                containersPaused, containersStopped, images, driver, labels
        );

        assertEquals(id, response.id());
        assertEquals(name, response.name());
        assertEquals(serverVersion, response.serverVersion());
        assertEquals(kernelVersion, response.kernelVersion());
        assertEquals(operatingSystem, response.operatingSystem());
        assertEquals(architecture, response.architecture());
        assertEquals(nCpu, response.nCpu());
        assertEquals(memTotal, response.memTotal());
        assertEquals(containers, response.containers());
        assertEquals(containersRunning, response.containersRunning());
        assertEquals(containersPaused, response.containersPaused());
        assertEquals(containersStopped, response.containersStopped());
        assertEquals(images, response.images());
        assertEquals(driver, response.driver());
        assertEquals(labels, response.labels());
    }

    @Test
    void systemInfoResponse_ShouldAllowNullValues() {
        SystemInfoResponse response = new SystemInfoResponse(
                null, null, null, null, null,
                null, null, null, null, null,
                null, null, null, null, null
        );

        assertNull(response.id());
        assertNull(response.name());
        assertNull(response.serverVersion());
        assertNull(response.kernelVersion());
        assertNull(response.operatingSystem());
        assertNull(response.architecture());
        assertNull(response.nCpu());
        assertNull(response.memTotal());
        assertNull(response.containers());
        assertNull(response.containersRunning());
        assertNull(response.containersPaused());
        assertNull(response.containersStopped());
        assertNull(response.images());
        assertNull(response.driver());
        assertNull(response.labels());
    }

    @Test
    void systemInfoResponse_ShouldSupportEquality() {
        List<String> labels = List.of("label1", "label2");
        SystemInfoResponse response1 = new SystemInfoResponse(
                "id1", "name1", "27.0.0", "6.8.0", "Ubuntu",
                "x86_64", 8, 16000000000L, 10, 3,
                0, 7, 5, "overlay2", labels
        );
        SystemInfoResponse response2 = new SystemInfoResponse(
                "id1", "name1", "27.0.0", "6.8.0", "Ubuntu",
                "x86_64", 8, 16000000000L, 10, 3,
                0, 7, 5, "overlay2", labels
        );
        SystemInfoResponse response3 = new SystemInfoResponse(
                "id2", "name1", "27.0.0", "6.8.0", "Ubuntu",
                "x86_64", 8, 16000000000L, 10, 3,
                0, 7, 5, "overlay2", labels
        );

        assertEquals(response1, response2);
        assertNotEquals(response1, response3);
        assertEquals(response1.hashCode(), response2.hashCode());
    }

    @Test
    void systemInfoResponse_ShouldHaveProperToString() {
        SystemInfoResponse response = new SystemInfoResponse(
                "id1", "name1", "27.0.0", "6.8.0", "Ubuntu",
                "x86_64", 8, 16000000000L, 10, 3,
                0, 7, 5, "overlay2", List.of("label1")
        );

        String toString = response.toString();

        assertNotNull(toString);
        assertTrue(toString.contains("id1"));
        assertTrue(toString.contains("name1"));
        assertTrue(toString.contains("27.0.0"));
        assertTrue(toString.contains("Ubuntu"));
        assertTrue(toString.contains("x86_64"));
    }
}
