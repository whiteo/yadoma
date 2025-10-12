package dev.whiteo.yadoma.security;

import com.fasterxml.jackson.databind.ObjectMapper;
import dev.whiteo.yadoma.dto.user.UserLoginRequest;
import jakarta.servlet.ReadListener;
import jakarta.servlet.ServletInputStream;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.junit.jupiter.api.Test;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.mockito.Mockito.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class JwtAuthenticationFilterTest {

    private static class SimpleServletInputStream extends ServletInputStream {
        private final InputStream delegate;

        SimpleServletInputStream(InputStream delegate) {
            this.delegate = delegate;
        }

        @Override
        public int read() throws IOException {
            return delegate.read();
        }

        @Override
        public boolean isFinished() {
            try {
                return delegate.available() == 0;
            } catch (IOException _) {
                return true;
            }
        }

        @Override
        public boolean isReady() {
            return true;
        }

        @Override
        public void setReadListener(ReadListener readListener) {}
    }

    @Test
    void attemptAuthentication_parsesRequestAndDelegatesToAuthenticationManager() throws Exception {
        JwtAuthenticationFilter filter = new JwtAuthenticationFilter();

        AuthenticationManager authManager = mock(AuthenticationManager.class);
        filter.setAuthenticationManager(authManager);

        UserLoginRequest reqObj = new UserLoginRequest("user@example.com", "passwd");
        String json = new ObjectMapper().writeValueAsString(reqObj);

        HttpServletRequest request = mock(HttpServletRequest.class);
        HttpServletResponse response = mock(HttpServletResponse.class);

        InputStream is = new ByteArrayInputStream(json.getBytes());
        ServletInputStream sis = new SimpleServletInputStream(is);

        when(request.getInputStream()).thenReturn(sis);

        Authentication expectedAuth = new UsernamePasswordAuthenticationToken(reqObj.email(), reqObj.password());
        when(authManager.authenticate(any(UsernamePasswordAuthenticationToken.class))).thenReturn(expectedAuth);

        Authentication result = filter.attemptAuthentication(request, response);

        assertNotNull(result);
        assertEquals(reqObj.email(), result.getPrincipal());
    }
}
