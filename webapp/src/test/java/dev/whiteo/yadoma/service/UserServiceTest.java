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

import java.util.ArrayList;
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

    @Test
    void getUserById_shouldReturnUser() {
        when(repository.getOrThrow("userId")).thenReturn(user);

        User result = userService.getUserById("userId");

        assertNotNull(result);
        assertEquals(user, result);
        verify(repository).getOrThrow("userId");
    }

    @Test
    void removeContainerFromUser_shouldRemoveAndSave() {
        user.setContainerIds(new ArrayList<>(Arrays.asList("container1", "container2")));

        userService.removeContainerFromUser("container1", user);

        assertEquals(1, user.getContainerIds().size());
        assertFalse(user.getContainerIds().contains("container1"));
        verify(repository).save(user);
    }

    @Test
    void addContainerToUser_shouldAddAndSave() {
        ArrayList<String> containerIds = new ArrayList<>(Arrays.asList("container1"));
        user.setContainerIds(containerIds);
        when(repository.save(any(User.class))).thenReturn(user);

        userService.addContainerToUser("container2", user);

        assertEquals(2, containerIds.size());
        assertTrue(containerIds.contains("container2"));
        verify(repository).save(user);
    }

    @Test
    void addContainerToUser_withNullList_shouldInitializeAndAdd() {
        user.setContainerIds(null);
        when(repository.save(any(User.class))).thenReturn(user);

        userService.addContainerToUser("container1", user);

        assertNotNull(user.getContainerIds());
        assertEquals(1, user.getContainerIds().size());
        assertTrue(user.getContainerIds().contains("container1"));
        verify(repository).save(user);
    }

    @Test
    void findAll_shouldReturnAllUsers_whenAdmin() {
        User admin = new User();
        admin.setId("adminId");
        admin.setRole(Role.ADMIN);

        User user2 = new User();
        user2.setId("user2Id");
        user2.setEmail("user2@example.com");

        when(repository.getOrThrow("adminId")).thenReturn(admin);
        when(repository.findAll()).thenReturn(Arrays.asList(admin, user, user2));
        when(mapper.toResponse(any(User.class))).thenReturn(userResponse);

        var result = userService.findAll("adminId");

        assertEquals(3, result.size());
        verify(repository).findAll();
    }

    @Test
    void findAll_shouldThrowException_whenNotAdmin() {
        user.setRole(Role.USER);
        when(repository.getOrThrow("userId")).thenReturn(user);

        assertThrows(BadCredentialsException.class, () -> userService.findAll("userId"));
        verify(repository).getOrThrow("userId");
    }

    @Test
    void remove_shouldDeleteUser_whenAdmin() {
        User admin = new User();
        admin.setId("adminId");
        admin.setRole(Role.ADMIN);

        when(repository.getOrThrow("adminId")).thenReturn(admin);
        when(repository.getOrThrow("userId")).thenReturn(user);

        assertDoesNotThrow(() -> userService.remove("userId", "adminId"));

        verify(repository).delete(user);
    }

    @Test
    void updatePassword_shouldUpdatePassword() {
        try (MockedStatic<PasswordUtil> passwordUtil = mockStatic(PasswordUtil.class)) {
            passwordUtil.when(() -> PasswordUtil.matches("oldPassword", "hashedPassword")).thenReturn(true);
            passwordUtil.when(() -> PasswordUtil.hash("newPassword")).thenReturn("newHashedPassword");

            when(repository.getOrThrow("userId")).thenReturn(user);

            var request = new dev.whiteo.yadoma.dto.user.UserUpdatePasswordRequest("oldPassword", "newPassword");

            assertDoesNotThrow(() -> userService.updatePassword("userId", request));

            verify(repository).save(user);
            passwordUtil.verify(() -> PasswordUtil.hash("newPassword"));
        }
    }

    @Test
    void validateUserAccess_shouldReturnUser_whenUserOwnsContainer() {
        user.setContainerIds(new ArrayList<>(Arrays.asList("container1")));
        when(repository.getOrThrow("userId")).thenReturn(user);

        User result = userService.validateUserAccess("container1", "userId");

        assertEquals(user, result);
    }

    @Test
    void validateUserAccess_shouldReturnUser_whenAdmin() {
        user.setRole(Role.ADMIN);
        user.setContainerIds(new ArrayList<>());
        when(repository.getOrThrow("userId")).thenReturn(user);

        User result = userService.validateUserAccess("anyContainer", "userId");

        assertEquals(user, result);
    }

    @Test
    void validateUserAccess_shouldThrowException_whenUserDoesNotOwnContainer() {
        user.setContainerIds(new ArrayList<>(Arrays.asList("container1")));
        when(repository.getOrThrow("userId")).thenReturn(user);

        assertThrows(BadCredentialsException.class, () ->
            userService.validateUserAccess("container2", "userId"));
    }
}