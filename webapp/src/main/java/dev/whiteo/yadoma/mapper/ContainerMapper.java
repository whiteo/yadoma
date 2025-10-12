package dev.whiteo.yadoma.mapper;

import container.v1.Container;
import dev.whiteo.yadoma.dto.container.ContainerResponse;
import org.mapstruct.Mapper;

/**
 * Mapper for converting User entities to UserResponse DTOs.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Mapper(componentModel = "spring")
public interface ContainerMapper extends AbstractMapper<Container, ContainerResponse> {

    /**
     * Converts a Container.GetContainerResponse to a ContainerResponse DTO.
     *
     * @param response the Container.GetContainerResponse to convert
     * @return the converted ContainerResponse DTO
     */
    ContainerResponse toResponseDTO(Container.GetContainerResponse response);
    /**
     * Converts a Container.GetContainerDetailsResponse to a ContainerResponse DTO.
     *
     * @param response the Container.GetContainerDetailsResponse to convert
     * @return the converted ContainerResponse DTO
     */
    ContainerResponse toResponseDTO(Container.GetContainerDetailsResponse response);
}