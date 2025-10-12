package dev.whiteo.yadoma.mapper;

import container.v1.Container;
import dev.whiteo.yadoma.dto.container.ContainerResponse;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

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
    @Mapping(target = "name", expression = "java(response.getNamesCount() > 0 ? response.getNames(0).replaceFirst(\"^/\", \"\") : \"\")")
    @Mapping(target = "createdAt", expression = "java(null)")
    ContainerResponse toResponseDTO(Container.GetContainerResponse response);

    /**
     * Converts a Container.GetContainerDetailsResponse to a ContainerResponse DTO.
     *
     * @param response the Container.GetContainerDetailsResponse to convert
     * @return the converted ContainerResponse DTO
     */
    @Mapping(target = "name", expression = "java(response.getName().replaceFirst(\"^/\", \"\"))")
    @Mapping(target = "createdAt", source = "created")
    ContainerResponse toResponseDTO(Container.GetContainerDetailsResponse response);
}