package dev.whiteo.yadoma.config;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.web.socket.server.standard.ServerEndpointExporter;

import static org.junit.jupiter.api.Assertions.assertNotNull;

class WebSocketConfigTest {

    private WebSocketConfig webSocketConfig;

    @BeforeEach
    void setUp() {
        webSocketConfig = new WebSocketConfig();
    }

    @Test
    void serverEndpointExporter_shouldCreateExporter() {
        ServerEndpointExporter exporter = webSocketConfig.serverEndpointExporter();

        assertNotNull(exporter);
    }
}
