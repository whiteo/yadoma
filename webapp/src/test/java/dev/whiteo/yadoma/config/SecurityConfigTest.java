package dev.whiteo.yadoma.config;

import dev.whiteo.yadoma.security.AuthInterceptor;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.junit.jupiter.api.Test;
import org.springframework.http.HttpStatus;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.web.authentication.logout.LogoutSuccessHandler;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.RETURNS_DEEP_STUBS;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;

class SecurityConfigTest {
    @Test
    void canInstantiate() {
        AuthInterceptor authInterceptor = mock(AuthInterceptor.class);
        assertDoesNotThrow(() -> new SecurityConfig(authInterceptor));
    }

    @Test
    void passwordEncoder_returnsBCryptPasswordEncoder() {
        AuthInterceptor authInterceptor = mock(AuthInterceptor.class);
        SecurityConfig config = new SecurityConfig(authInterceptor);
        PasswordEncoder encoder = config.passwordEncoder();
        assertNotNull(encoder);
        String raw = "pass";
        String encoded = encoder.encode(raw);
        assertTrue(encoder.matches(raw, encoded));
    }

    @Test
    void logoutSuccessHandler_setsStatusOk() throws Exception {
        AuthInterceptor authInterceptor = mock(AuthInterceptor.class);
        SecurityConfig config = new SecurityConfig(authInterceptor);
        LogoutSuccessHandler handler = config.logoutSuccessHandler();
        var request = mock(HttpServletRequest.class);
        var response = mock(HttpServletResponse.class);
        handler.onLogoutSuccess(request, response, null);
        verify(response).setStatus(HttpStatus.OK.value());
    }

    @Test
    void filterChain_buildsSecurityFilterChain() throws Exception {
        AuthInterceptor authInterceptor = mock(AuthInterceptor.class);
        SecurityConfig config = new SecurityConfig(authInterceptor);
        HttpSecurity http = mock(HttpSecurity.class, RETURNS_DEEP_STUBS);
        assertNotNull(config.filterChain(http));
    }
}
