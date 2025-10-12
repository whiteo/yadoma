package dev.whiteo.yadoma.config;

import jakarta.websocket.server.ServerEndpointConfig;
import org.springframework.beans.BeansException;
import org.springframework.beans.factory.BeanFactory;
import org.springframework.context.ApplicationContext;
import org.springframework.context.ApplicationContextAware;
import org.springframework.stereotype.Component;

/**
 * Configurator for Jakarta WebSocket endpoints to enable Spring dependency injection.
 * This allows @ServerEndpoint classes to access Spring beans.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Component
public class SpringWebSocketConfigurator extends ServerEndpointConfig.Configurator implements ApplicationContextAware {

    private static volatile BeanFactory context;

    @Override
    public <T> T getEndpointInstance(Class<T> endpointClass) throws InstantiationException {
        return context.getBean(endpointClass);
    }

    @Override
    public void setApplicationContext(ApplicationContext applicationContext) throws BeansException {
        SpringWebSocketConfigurator.context = applicationContext;
    }
}
