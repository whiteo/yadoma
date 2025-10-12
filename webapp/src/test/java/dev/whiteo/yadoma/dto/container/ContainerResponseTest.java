package dev.whiteo.yadoma.dto.container;

import org.junit.jupiter.api.Test;

import java.time.LocalDateTime;

import static org.junit.jupiter.api.Assertions.*;

class ContainerResponseTest {

    @Test
    void containerResponse_ShouldCreateWithAllFields() {
        // Given
        String id = "container123";
        String name = "test-container";
        LocalDateTime createdAt = LocalDateTime.now();
        String status = "running";
        String state = "active";

        // When
        ContainerResponse response = new ContainerResponse(id, name, createdAt, status, state);

        // Then
        assertEquals(id, response.id());
        assertEquals(name, response.name());
        assertEquals(createdAt, response.createdAt());
        assertEquals(status, response.status());
        assertEquals(state, response.state());
    }

    @Test
    void containerResponse_ShouldAllowNullValues() {
        // When
        ContainerResponse response = new ContainerResponse(null, null, null, null, null);

        // Then
        assertNull(response.id());
        assertNull(response.name());
        assertNull(response.createdAt());
        assertNull(response.status());
        assertNull(response.state());
    }

    @Test
    void containerResponse_ShouldSupportEquality() {
        // Given
        LocalDateTime now = LocalDateTime.now();
        ContainerResponse response1 = new ContainerResponse("id1", "name1", now, "running", "active");
        ContainerResponse response2 = new ContainerResponse("id1", "name1", now, "running", "active");
        ContainerResponse response3 = new ContainerResponse("id2", "name1", now, "running", "active");

        // Then
        assertEquals(response1, response2);
        assertNotEquals(response1, response3);
        assertEquals(response1.hashCode(), response2.hashCode());
    }

    @Test
    void containerResponse_ShouldHaveProperToString() {
        // Given
        LocalDateTime now = LocalDateTime.now();
        ContainerResponse response = new ContainerResponse("id1", "name1", now, "running", "active");

        // When
        String toString = response.toString();

        // Then
        assertNotNull(toString);
        assertTrue(toString.contains("id1"));
        assertTrue(toString.contains("name1"));
        assertTrue(toString.contains("running"));
        assertTrue(toString.contains("active"));
    }
}
