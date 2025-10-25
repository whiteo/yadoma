package dev.whiteo.yadoma.config;

import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.boot.autoconfigure.condition.ConditionalOnWebApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.socket.server.standard.ServerEndpointExporter;

/**
 * WebSocket configuration for Jakarta WebSocket endpoints.
 * Registers @ServerEndpoint annotated classes with the container.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Configuration
@ConditionalOnProperty(value = "app.websocket.enabled", havingValue = "true", matchIfMissing = true)
public class WebSocketConfig {

    @Bean
    @ConditionalOnWebApplication
    public ServerEndpointExporter serverEndpointExporter() {
        return new ServerEndpointExporter();
    }
}
