package dev.whiteo.yadoma.service;

import dev.whiteo.yadoma.domain.Role;
import dev.whiteo.yadoma.domain.User;
import dev.whiteo.yadoma.dto.token.TokenResponse;
import dev.whiteo.yadoma.dto.token.TokenValidationResponse;
import dev.whiteo.yadoma.dto.user.UserCreateRequest;
import dev.whiteo.yadoma.dto.user.UserLoginRequest;
import dev.whiteo.yadoma.dto.user.UserResponse;
import dev.whiteo.yadoma.mapper.UserMapper;
import dev.whiteo.yadoma.repository.UserRepository;
import dev.whiteo.yadoma.security.TokenInteract;
import dev.whiteo.yadoma.security.UserDetailsImpl;
import dev.whiteo.yadoma.util.PasswordUtil;
import jakarta.servlet.http.HttpServletRequest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockedStatic;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.core.userdetails.UserDetails;

import java.util.Arrays;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.mockStatic;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class UserServiceTest {

    @Mock
    private TokenInteract tokenInteract;

    @Mock
    private UserRepository repository;

    @Mock
    private UserMapper mapper;

    @InjectMocks
    private UserService userService;

    private User user;
    private UserLoginRequest loginRequest;
    private UserResponse userResponse;
    private TokenResponse tokenResponse;

    @BeforeEach
    void setUp() {
        user = new User();
        user.setId("userId");
        user.setEmail("test@example.com");
        user.setPasswordHash("hashedPassword");
        user.setRole(Role.USER);

        loginRequest = new UserLoginRequest("test@example.com", "password");

        userResponse = new UserResponse("userId","test@example.com", Role.USER);

        tokenResponse = new TokenResponse("token");
    }

    @Test
    void getToken_WithValidCredentials_ShouldReturnToken() {
        try (MockedStatic<PasswordUtil> passwordUtil = mockStatic(PasswordUtil.class)) {
            passwordUtil.when(() -> PasswordUtil.matches("password", "hashedPassword")).thenReturn(true);

            when(repository.findByEmailIgnoreCase("test@example.com")).thenReturn(java.util.Optional.of(user));
            when(repository.getOrThrow("userId")).thenReturn(user); // Add mock for loadUserByUsername
            when(tokenInteract.generateToken(any(UserDetails.class))).thenReturn("token");

            TokenResponse result = userService.getToken(loginRequest);

            assertNotNull(result);
            assertEquals("token", result.token());
            passwordUtil.verify(() -> PasswordUtil.matches("password", "hashedPassword"));
            verify(tokenInteract).generateToken(any(UserDetails.class));
        }
    }

    @Test
    void getToken_WithInvalidPassword_ShouldThrowException() {
        try (MockedStatic<PasswordUtil> passwordUtil = mockStatic(PasswordUtil.class)) {
            passwordUtil.when(() -> PasswordUtil.matches("password", "hashedPassword")).thenReturn(false);

            when(repository.findByEmailIgnoreCase("test@example.com")).thenReturn(java.util.Optional.of(user));

            assertThrows(BadCredentialsException.class, () -> userService.getToken(loginRequest));
            passwordUtil.verify(() -> PasswordUtil.matches("password", "hashedPassword"));
        }
    }

    @Test
    void validateToken_shouldReturnValidationResult() {
        HttpServletRequest request = mock(HttpServletRequest.class);
        when(tokenInteract.getToken(request)).thenReturn("token");
        when(tokenInteract.validateToken("token")).thenReturn(true);

        TokenValidationResponse result = userService.validateToken(request);

        assertTrue(result.valid());
        verify(tokenInteract).validateToken("token");
    }

    @Test
    void loadUserByUsername_shouldReturnUserDetails() {
        when(repository.getOrThrow("userId")).thenReturn(user);

        UserDetails result = userService.loadUserByUsername("userId");

        assertNotNull(result);
        assertTrue(result instanceof UserDetailsImpl);
        assertEquals("userId", result.getUsername());
    }

    @Test
    void loadUserByUsername_withNonExistentUser_shouldThrowException() {
        when(repository.getOrThrow("unknown")).thenThrow(new RuntimeException("User not found"));

        assertThrows(RuntimeException.class, () -> userService.loadUserByUsername("unknown"));
    }

    @Test
    void me_shouldReturnUserResponse() {
        when(repository.getOrThrow("userId")).thenReturn(user);
        when(mapper.toResponse(user)).thenReturn(userResponse);

        UserResponse result = userService.me("userId");

        assertEquals(userResponse, result);
        verify(repository).getOrThrow("userId");
        verify(mapper).toResponse(user);
    }

    @Test
    void create_shouldCreateUser() {
        UserCreateRequest request = new UserCreateRequest("new@example.com", "password");

        try (MockedStatic<PasswordUtil> passwordUtil = mockStatic(PasswordUtil.class)) {
            passwordUtil.when(() -> PasswordUtil.hash("password")).thenReturn("hashedPassword");
            when(repository.findByEmailIgnoreCase("new@example.com")).thenReturn(java.util.Optional.empty());
            when(mapper.toEntity("new@example.com", "hashedPassword")).thenReturn(user);

            assertDoesNotThrow(() -> userService.create(request));

            passwordUtil.verify(() -> PasswordUtil.hash("password"));
            verify(repository).save(user);
        }
    }

    @Test
    void isContainerIdContains_shouldReturnTrue_whenUserOwnsContainer() {
        user.setContainerIds(Arrays.asList("container1", "container2"));
        when(repository.getOrThrow("userId")).thenReturn(user);

        Boolean result = userService.isContainerIdContains("container1", "userId");

        assertTrue(result);
    }

    @Test
    void isContainerIdContains_shouldReturnFalse_whenUserDoesNotOwnContainer() {
        user.setContainerIds(Arrays.asList("container1", "container2"));
        when(repository.getOrThrow("userId")).thenReturn(user);

        Boolean result = userService.isContainerIdContains("container3", "userId");

        assertFalse(result);
    }

    @Test
    void isContainerIdContains_shouldReturnTrue_whenUserIsAdmin() {
        user.setRole(Role.ADMIN);
        user.setContainerIds(Arrays.asList("container1"));
        when(repository.getOrThrow("userId")).thenReturn(user);

        Boolean result = userService.isContainerIdContains("container3", "userId");

        assertTrue(result);
    }

    @Test
    void validateToken_invalidToken_shouldReturnFalse() {
        HttpServletRequest request = mock(HttpServletRequest.class);
        when(tokenInteract.getToken(request)).thenReturn("invalidtoken");
        when(tokenInteract.validateToken("invalidtoken")).thenReturn(false);

        TokenValidationResponse result = userService.validateToken(request);

        assertFalse(result.valid());
        verify(tokenInteract).validateToken("invalidtoken");
    }
}