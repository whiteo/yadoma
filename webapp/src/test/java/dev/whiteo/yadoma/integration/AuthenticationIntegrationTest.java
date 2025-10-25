package dev.whiteo.yadoma.integration;

import com.fasterxml.jackson.databind.ObjectMapper;
import dev.whiteo.yadoma.domain.User;
import dev.whiteo.yadoma.dto.token.TokenResponse;
import dev.whiteo.yadoma.dto.user.UserCreateRequest;
import dev.whiteo.yadoma.dto.user.UserLoginRequest;
import dev.whiteo.yadoma.repository.UserRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.MediaType;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.MvcResult;
import org.testcontainers.containers.MongoDBContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import static org.junit.jupiter.api.Assertions.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

/**
 * Integration tests for authentication endpoints using Testcontainers.
 * Tests the complete authentication flow including user registration and login.
 */
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@AutoConfigureMockMvc
@Testcontainers
@ActiveProfiles("test")
class AuthenticationIntegrationTest {

    @Container
    static MongoDBContainer mongoDBContainer = new MongoDBContainer("mongo:8.0")
            .withExposedPorts(27017);

    @DynamicPropertySource
    static void setProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.data.mongodb.uri", mongoDBContainer::getReplicaSetUrl);
    }

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private ObjectMapper objectMapper;

    @Autowired
    private UserRepository userRepository;

    @BeforeEach
    void setUp() {
        userRepository.deleteAll();
    }

    @Test
    void shouldRegisterNewUser() throws Exception {
        UserCreateRequest request = new UserCreateRequest(
                "test@example.com",
                "SecurePass123!"
        );

        mockMvc.perform(post("/api/v1/user/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk());

        User savedUser = userRepository.findByEmailIgnoreCase("test@example.com").orElse(null);
        assertNotNull(savedUser);
        assertEquals("test@example.com", savedUser.getEmail());
    }

    @Test
    void shouldLoginWithValidCredentials() throws Exception {
        
        UserCreateRequest registerRequest = new UserCreateRequest(
                "login@example.com",
                "SecurePass123!"
        );

        mockMvc.perform(post("/api/v1/user/register")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(registerRequest)));

        
        UserLoginRequest loginRequest = new UserLoginRequest("login@example.com", "SecurePass123!");

        MvcResult result = mockMvc.perform(post("/api/v1/authenticate/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginRequest)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.token").exists())
                .andReturn();

        String responseContent = result.getResponse().getContentAsString();
        TokenResponse tokenResponse = objectMapper.readValue(responseContent, TokenResponse.class);

        assertNotNull(tokenResponse.token());
        assertFalse(tokenResponse.token().isEmpty());
    }

    @Test
    void shouldRejectLoginWithInvalidPassword() throws Exception {
        
        UserCreateRequest registerRequest = new UserCreateRequest(
                "invalid@example.com",
                "CorrectPass123!"
        );

        mockMvc.perform(post("/api/v1/user/register")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(registerRequest)));

        
        UserLoginRequest loginRequest = new UserLoginRequest("invalid@example.com", "WrongPassword");

        mockMvc.perform(post("/api/v1/authenticate/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginRequest)))
                .andExpect(status().isUnauthorized());
    }

    @Test
    void shouldRejectLoginWithNonExistentUser() throws Exception {
        UserLoginRequest loginRequest = new UserLoginRequest("nonexistent@example.com", "anypassword");

        mockMvc.perform(post("/api/v1/authenticate/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginRequest)))
                .andExpect(status().isUnauthorized());
    }

    @Test
    void shouldRejectDuplicateEmail() throws Exception {
        UserCreateRequest request1 = new UserCreateRequest(
                "duplicate@example.com",
                "SecurePass123!"
        );

        
        mockMvc.perform(post("/api/v1/user/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request1)))
                .andExpect(status().isOk());

        UserCreateRequest request2 = new UserCreateRequest(
                "duplicate@example.com",
                "DifferentPass456!"
        );

        
        mockMvc.perform(post("/api/v1/user/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request2)))
                .andExpect(status().is4xxClientError());
    }

    @Test
    void shouldValidateTokenAfterLogin() throws Exception {
        
        UserCreateRequest registerRequest = new UserCreateRequest(
                "token@example.com",
                "SecurePass123!"
        );

        mockMvc.perform(post("/api/v1/user/register")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(registerRequest)));

        UserLoginRequest loginRequest = new UserLoginRequest("token@example.com", "SecurePass123!");

        MvcResult loginResult = mockMvc.perform(post("/api/v1/authenticate/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginRequest)))
                .andExpect(status().isOk())
                .andReturn();

        String responseContent = loginResult.getResponse().getContentAsString();
        TokenResponse tokenResponse = objectMapper.readValue(responseContent, TokenResponse.class);

        
        mockMvc.perform(post("/api/v1/authenticate/validate")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("{\"token\":\"" + tokenResponse.token() + "\"}"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.valid").value(true))
                .andExpect(jsonPath("$.userId").exists());
    }
}
