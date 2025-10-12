package dev.whiteo.yadoma.mapper;

import dev.whiteo.yadoma.domain.User;
import dev.whiteo.yadoma.dto.user.UserResponse;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

/**
 * Mapper for converting User entities to UserResponse DTO.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Mapper(componentModel = "spring")
public interface UserMapper extends AbstractMapper<User, UserResponse> {

    /**
     * Converts email and password hash to a User entity.
     *
     * @param email        user's email
     * @param passwordHash hashed password
     * @return User entity
     */
    @Mapping(target = "id", ignore = true)
    @Mapping(target = "role", ignore = true)
    @Mapping(target = "email", source = "email")
    @Mapping(target = "modifyDate", ignore = true)
    @Mapping(target = "creationDate", ignore = true)
    @Mapping(target = "containerIds", ignore = true)
    @Mapping(target = "passwordHash", source = "passwordHash")
    User toEntity(String email, String passwordHash);
}