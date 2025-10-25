package dev.whiteo.yadoma.integration;

import com.fasterxml.jackson.databind.ObjectMapper;
import dev.whiteo.yadoma.domain.User;
import dev.whiteo.yadoma.dto.token.TokenResponse;
import dev.whiteo.yadoma.dto.user.UserCreateRequest;
import dev.whiteo.yadoma.dto.user.UserLoginRequest;
import dev.whiteo.yadoma.repository.UserRepository;
import org.junit.jupiter.api.BeforeEach;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.annotation.Import;
import org.springframework.http.MediaType;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.MvcResult;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;

/**
 * Integration tests for user management endpoints using Testcontainers.
 */
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@AutoConfigureMockMvc
@ActiveProfiles("test")
@Import(dev.whiteo.yadoma.config.TestSecurityConfig.class)
class UserManagementIntegrationTest {

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private ObjectMapper objectMapper;

    @Autowired
    private UserRepository userRepository;

    private String authToken;
    private String userId;

    @BeforeEach
    void setUp() throws Exception {
        userRepository.deleteAll();

        
        UserCreateRequest registerRequest = new UserCreateRequest(
                "test@example.com",
                "SecurePass123!"
        );

        mockMvc.perform(post("/api/v1/user/create")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(registerRequest)));

        UserLoginRequest loginRequest = new UserLoginRequest("test@example.com", "SecurePass123!");

        MvcResult loginResult = mockMvc.perform(post("/api/v1/authenticate")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginRequest)))
                .andReturn();

        String responseContent = loginResult.getResponse().getContentAsString();
        TokenResponse tokenResponse = objectMapper.readValue(responseContent, TokenResponse.class);
        authToken = tokenResponse.token();

        User user = userRepository.findByEmailIgnoreCase("test@example.com").orElseThrow();
        userId = user.getId();
    }
}
