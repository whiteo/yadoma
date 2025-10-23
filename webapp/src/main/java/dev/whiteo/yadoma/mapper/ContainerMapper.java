package dev.whiteo.yadoma.mapper;

import container.v1.Container;
import dev.whiteo.yadoma.dto.container.ContainerResponse;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.time.LocalDateTime;
import java.time.ZonedDateTime;
import java.time.format.DateTimeFormatter;

/**
 * Mapper for converting Container proto messages to ContainerResponse DTOs.
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
    @Mapping(target = "createdAt", expression = "java(parseDateTime(response.getCreated()))")
    ContainerResponse toResponseDTO(Container.GetContainerDetailsResponse response);

    /**
     * Parses ISO 8601 datetime string to LocalDateTime.
     *
     * @param dateTimeStr the datetime string (e.g., "2023-01-01T00:00:00Z")
     * @return LocalDateTime or null if parsing fails
     */
    default LocalDateTime parseDateTime(String dateTimeStr) {
        if (dateTimeStr == null || dateTimeStr.isEmpty()) {
            return null;
        }
        try {
            return ZonedDateTime.parse(dateTimeStr, DateTimeFormatter.ISO_DATE_TIME).toLocalDateTime();
        } catch (Exception e) {
            return null;
        }
    }
}