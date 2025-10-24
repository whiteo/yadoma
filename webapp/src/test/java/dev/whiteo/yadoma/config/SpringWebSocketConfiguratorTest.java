package dev.whiteo.yadoma.config;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.context.ApplicationContext;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class SpringWebSocketConfiguratorTest {

    private SpringWebSocketConfigurator configurator;
    private ApplicationContext applicationContext;

    @BeforeEach
    void setUp() {
        configurator = new SpringWebSocketConfigurator();
        applicationContext = mock(ApplicationContext.class);
    }

    @Test
    void setApplicationContext_shouldSetContext() {
        configurator.setApplicationContext(applicationContext);

        assertNotNull(configurator);
    }

    @Test
    void getEndpointInstance_shouldReturnBeanFromContext() throws InstantiationException {
        TestEndpoint testEndpoint = new TestEndpoint();
        when(applicationContext.getBean(TestEndpoint.class)).thenReturn(testEndpoint);

        configurator.setApplicationContext(applicationContext);
        TestEndpoint result = configurator.getEndpointInstance(TestEndpoint.class);

        assertEquals(testEndpoint, result);
    }

    static class TestEndpoint {
    }
}
