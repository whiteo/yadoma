package dev.whiteo.yadoma.controller;

import dev.whiteo.yadoma.dto.token.TokenResponse;
import dev.whiteo.yadoma.dto.token.TokenValidationResponse;
import dev.whiteo.yadoma.dto.user.UserLoginRequest;
import dev.whiteo.yadoma.service.UserService;
import jakarta.servlet.http.HttpServletRequest;
import org.junit.jupiter.api.Test;
import org.springframework.http.HttpStatus;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class AuthenticateRestControllerTest {

    @Test
    void authenticate_shouldReturnLoginResponse() {
        UserService service = mock(UserService.class);
        AuthenticateRestController controller = new AuthenticateRestController(service);

        UserLoginRequest request = new UserLoginRequest("test@example.com", "pass");
        TokenResponse response = new TokenResponse("token");

        when(service.getToken(request)).thenReturn(response);

        var result = controller.authenticate(request);

        assertNotNull(result);
        assertEquals(HttpStatus.OK, result.getStatusCode());
        assertNotNull(result.getBody());
        assertEquals(response, result.getBody());
    }

    @Test
    void validate_shouldReturnTokenValidationResponse() {
        UserService service = mock(UserService.class);
        AuthenticateRestController controller = new AuthenticateRestController(service);

        HttpServletRequest request = mock(HttpServletRequest.class);
        TokenValidationResponse validationResponse = new TokenValidationResponse(true);
        when(service.validateToken(request)).thenReturn(validationResponse);

        var response = controller.validate(request);

        assertNotNull(response);
        assertEquals(HttpStatus.OK, response.getStatusCode());
        assertNotNull(response.getBody());
        assertTrue(response.getBody().valid());
    }
}
