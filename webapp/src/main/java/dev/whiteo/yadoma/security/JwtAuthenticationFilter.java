package dev.whiteo.yadoma.security;

import com.fasterxml.jackson.databind.ObjectMapper;
import dev.whiteo.yadoma.dto.user.UserLoginRequest;
import dev.whiteo.yadoma.exception.AppRuntimeException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;

import java.io.IOException;

/**
 * Authentication filter for processing JWT login requests.
 * Reads credentials from request and authenticates using AuthenticationManager.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@RequiredArgsConstructor
public class JwtAuthenticationFilter extends UsernamePasswordAuthenticationFilter {

    /**
     * Attempts authentication by reading UserLoginRequest from the request.
     * @param request HTTP servlet request
     * @param response HTTP servlet response
     * @return Authentication object
     * @throws org.springframework.security.core.AuthenticationException if authentication fails
     */
    @Override
    public Authentication attemptAuthentication(HttpServletRequest request, HttpServletResponse response)
            throws AuthenticationException {
        try {
            UserLoginRequest loginRequest = new ObjectMapper()
                    .readValue(request.getInputStream(), UserLoginRequest.class);
            UsernamePasswordAuthenticationToken authToken = new UsernamePasswordAuthenticationToken(
                    loginRequest.email(), loginRequest.password());
            return getAuthenticationManager().authenticate(authToken);
        } catch (IOException e) {
            throw new AppRuntimeException(e);
        }
    }
}