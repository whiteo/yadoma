package dev.whiteo.yadoma.controller;

import dev.whiteo.yadoma.domain.Role;
import dev.whiteo.yadoma.dto.user.UserCreateRequest;
import dev.whiteo.yadoma.dto.user.UserResponse;
import dev.whiteo.yadoma.dto.user.UserUpdatePasswordRequest;
import dev.whiteo.yadoma.security.AuthInterceptor;
import dev.whiteo.yadoma.service.UserService;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;
import org.springframework.http.ResponseEntity;

import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.mockito.Mockito.doNothing;
import static org.mockito.Mockito.when;

class UserRestControllerTest {
    @Mock
    private UserService userService;
    @Mock
    private AuthInterceptor authInterceptor;
    @InjectMocks
    private UserRestController userRestController;

    @BeforeEach
    void setUp() {
        MockitoAnnotations.openMocks(this);
    }

    @Test
    void me_success() {
        String userId = "user123";
        UserResponse response = new UserResponse("test@email.com", Role.USER);
        when(authInterceptor.getUserId()).thenReturn(userId);
        when(userService.me(userId)).thenReturn(response);
        ResponseEntity<UserResponse> result = userRestController.me();
        assertEquals(org.springframework.http.HttpStatus.OK, result.getStatusCode());
        assertEquals(response, result.getBody());
    }

    @Test
    void create_success() {
        UserCreateRequest req = new UserCreateRequest("test@email.com", "password123");
        doNothing().when(userService).create(req);
        ResponseEntity<Void> result = userRestController.create(req);
        assertEquals(org.springframework.http.HttpStatus.CREATED, result.getStatusCode());
        assertNull(result.getBody());
    }

    @Test
    void updatePassword_success() {
        String userId = "user123";
        UserUpdatePasswordRequest req = new UserUpdatePasswordRequest("newPassword123", "oldPassword123");
        when(authInterceptor.getUserId()).thenReturn(userId);
        doNothing().when(userService).updatePassword(userId, req);
        ResponseEntity<Void> result = userRestController.updatePassword(req);
        assertEquals(org.springframework.http.HttpStatus.OK, result.getStatusCode());
        assertNull(result.getBody());
    }

    @Test
    void delete_success() {
        String userId = "user123";
        String deleteId = "deleteUser123";
        when(authInterceptor.getUserId()).thenReturn(userId);
        doNothing().when(userService).remove(deleteId, userId);
        ResponseEntity<Void> result = userRestController.delete(deleteId);
        assertEquals(org.springframework.http.HttpStatus.NO_CONTENT, result.getStatusCode());
        assertNull(result.getBody());
    }

    @Test
    void findAll_success() {
        String userId = "user123";
        UserResponse response1 = new UserResponse("test1@email.com", Role.USER);
        UserResponse response2 = new UserResponse("test2@email.com", Role.ADMIN);
        List<UserResponse> responses = List.of(response1, response2);

        when(authInterceptor.getUserId()).thenReturn(userId);
        when(userService.findAll(userId)).thenReturn(responses);

        ResponseEntity<List<UserResponse>> result = userRestController.findAll();

        assertEquals(org.springframework.http.HttpStatus.OK, result.getStatusCode());
        assertEquals(responses, result.getBody());
    }
}
