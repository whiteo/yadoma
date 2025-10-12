package dev.whiteo.yadoma.mapper;

import container.v1.Container;
import dev.whiteo.yadoma.dto.container.ContainerResponse;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mapstruct.factory.Mappers;

import static org.junit.jupiter.api.Assertions.*;

class ContainerMapperTest {

    private ContainerMapper containerMapper;

    @BeforeEach
    void setUp() {
        containerMapper = Mappers.getMapper(ContainerMapper.class);
    }

    @Test
    void toResponseDTO_ShouldMapGetContainerResponse() {
        Container.GetContainerResponse protoResponse = Container.GetContainerResponse.newBuilder()
                .setId("container123")
                .addNames("test-container")
                .setImage("nginx:latest")
                .setState("running")
                .setStatus("Up 5 minutes")
                .addPorts("80/tcp")
                .build();

        ContainerResponse result = containerMapper.toResponseDTO(protoResponse);

        assertNotNull(result);
        assertEquals("container123", result.id());
        assertEquals("running", result.state());
        assertEquals("Up 5 minutes", result.status());
    }

    @Test
    void toResponseDTO_ShouldMapGetContainerDetailsResponse() {
        Container.GetContainerDetailsResponse detailsResponse = Container.GetContainerDetailsResponse.newBuilder()
                .setId("container123")
                .setName("test-container")
                .setImage("nginx:latest")
                .setState("running")
                .setStatus("Up 5 minutes")
                .setCreated("2023-01-01T00:00:00Z")
                .build();

        ContainerResponse result = containerMapper.toResponseDTO(detailsResponse);

        assertNotNull(result);
        assertEquals("container123", result.id());
        assertEquals("test-container", result.name());
        assertEquals("running", result.state());
        assertEquals("Up 5 minutes", result.status());
    }

    @Test
    void toResponseDTO_ShouldHandleEmptyResponse() {
        Container.GetContainerResponse emptyResponse = Container.GetContainerResponse.newBuilder().build();

        ContainerResponse result = containerMapper.toResponseDTO(emptyResponse);

        assertNotNull(result);
    }
}
