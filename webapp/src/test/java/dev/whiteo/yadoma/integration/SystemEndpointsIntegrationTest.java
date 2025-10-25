package dev.whiteo.yadoma.integration;

import com.fasterxml.jackson.databind.ObjectMapper;
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

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

/**
 * Integration tests for system monitoring and version endpoints using Testcontainers.
 */
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@AutoConfigureMockMvc
@Testcontainers
@ActiveProfiles("test")
class SystemEndpointsIntegrationTest {

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

    private String authToken;

    @BeforeEach
    void setUp() throws Exception {
        userRepository.deleteAll();

        
        UserCreateRequest registerRequest = new UserCreateRequest(
                "sys@example.com",
                "SecurePass123!"
        );

        mockMvc.perform(post("/api/v1/user/register")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(registerRequest)));

        UserLoginRequest loginRequest = new UserLoginRequest("sys@example.com", "SecurePass123!");

        MvcResult loginResult = mockMvc.perform(post("/api/v1/authenticate/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(loginRequest)))
                .andReturn();

        String responseContent = loginResult.getResponse().getContentAsString();
        TokenResponse tokenResponse = objectMapper.readValue(responseContent, TokenResponse.class);
        authToken = tokenResponse.token();
    }

    @Test
    void shouldGetSystemInfo() throws Exception {
        mockMvc.perform(get("/api/v1/system/info")
                        .header("Authorization", "Bearer " + authToken))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.os").exists())
                .andExpect(jsonPath("$.architecture").exists())
                .andExpect(jsonPath("$.cpuCores").exists())
                .andExpect(jsonPath("$.totalMemory").exists())
                .andExpect(jsonPath("$.hostname").exists());
    }

    @Test
    void shouldGetDiskUsage() throws Exception {
        mockMvc.perform(get("/api/v1/system/disk-usage")
                        .header("Authorization", "Bearer " + authToken))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.containers").exists())
                .andExpect(jsonPath("$.images").exists())
                .andExpect(jsonPath("$.volumes").exists());
    }

    @Test
    void shouldRejectUnauthenticatedSystemInfoRequest() throws Exception {
        mockMvc.perform(get("/api/v1/system/info"))
                .andExpect(status().isUnauthorized());
    }

    @Test
    void shouldRejectUnauthenticatedDiskUsageRequest() throws Exception {
        mockMvc.perform(get("/api/v1/system/disk-usage"))
                .andExpect(status().isUnauthorized());
    }

    @Test
    void shouldGetVersion() throws Exception {
        
        mockMvc.perform(get("/api/v1/version"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.version").exists());
    }

    @Test
    void shouldAccessActuatorHealth() throws Exception {
        
        mockMvc.perform(get("/actuator/health"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.status").exists());
    }

    @Test
    void shouldAccessApiDocs() throws Exception {
        mockMvc.perform(get("/v3/api-docs"))
                .andExpect(status().isOk())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON));
    }
}
