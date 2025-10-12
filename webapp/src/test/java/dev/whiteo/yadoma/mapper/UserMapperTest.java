package dev.whiteo.yadoma.mapper;

import dev.whiteo.yadoma.domain.Role;
import dev.whiteo.yadoma.domain.User;
import dev.whiteo.yadoma.dto.user.UserResponse;
import org.junit.jupiter.api.Test;
import org.mapstruct.factory.Mappers;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;

class UserMapperTest {
    private final UserMapper userMapper = Mappers.getMapper(UserMapper.class);

    @Test
    void toResponse_shouldMapUserToUserResponse() {
        User user = new User();
        user.setEmail("test@example.com");
        user.setRole(Role.USER);

        UserResponse response = userMapper.toResponse(user);

        assertNotNull(response);
        assertEquals("test@example.com", response.email());
        assertEquals(Role.USER, response.role());
    }

    @Test
    void toEntity_shouldMapEmailAndPasswordToUser() {
        String email = "test@example.com";
        String passwordHash = "hashedPassword123";

        User user = userMapper.toEntity(email, passwordHash);

        assertNotNull(user);
        assertEquals(email, user.getEmail());
        assertEquals(passwordHash, user.getPasswordHash());
        assertEquals(Role.USER, user.getRole());
    }
}
