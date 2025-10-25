package dev.whiteo.yadoma.config;

import dev.whiteo.yadoma.security.AuthInterceptor;
import dev.whiteo.yadoma.security.TokenInteract;
import jakarta.servlet.http.HttpServletRequest;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Primary;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.web.SecurityFilterChain;

@TestConfiguration
public class TestSecurityConfig {

    @Bean
    @Primary
    public SecurityFilterChain testSecurityFilterChain(HttpSecurity http) throws Exception {
        http
                .csrf(AbstractHttpConfigurer::disable)
                .authorizeHttpRequests(auth -> auth.anyRequest().permitAll()
                );
        return http.build();
    }

    @Primary
    public AuthInterceptor authInterceptor() {
        return new AuthInterceptor(new DummyTokenInteract());
    }

    public static class DummyTokenInteract extends TokenInteract {
        public DummyTokenInteract() {
            super();
        }

        @Override
        public String getToken(HttpServletRequest request) {
            return "fake-token";
        }

        @Override
        public boolean validateToken(String token) {
            return true;
        }

        @Override
        public String getUserId(String token) {
            return "test-user";
        }
    }
}