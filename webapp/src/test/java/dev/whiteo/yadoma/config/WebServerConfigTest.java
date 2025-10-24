package dev.whiteo.yadoma.config;

import io.undertow.servlet.api.DeploymentInfo;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.boot.web.embedded.undertow.UndertowServletWebServerFactory;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;

class WebServerConfigTest {

    private WebServerConfig webServerConfig;

    @BeforeEach
    void setUp() {
        webServerConfig = new WebServerConfig();
    }

    @Test
    void customize_shouldAddDeploymentInfoCustomizer() {
        UndertowServletWebServerFactory factory = mock(UndertowServletWebServerFactory.class);

        assertDoesNotThrow(() -> webServerConfig.customize(factory));

        verify(factory).addDeploymentInfoCustomizers(any());
    }

    @Test
    void deploymentInfoCustomizer_shouldConfigureWebSocket() {
        UndertowServletWebServerFactory factory = new UndertowServletWebServerFactory();
        DeploymentInfo deploymentInfo = new DeploymentInfo();

        webServerConfig.customize(factory);

        assertDoesNotThrow(() -> factory.getDeploymentInfoCustomizers().forEach(c -> c.customize(deploymentInfo)));
    }
}
