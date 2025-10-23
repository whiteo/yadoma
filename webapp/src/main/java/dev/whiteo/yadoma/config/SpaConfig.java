package dev.whiteo.yadoma.config;

import org.springframework.context.annotation.Configuration;
import org.springframework.core.io.ClassPathResource;
import org.springframework.core.io.Resource;
import org.springframework.web.servlet.config.annotation.ResourceHandlerRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;
import org.springframework.web.servlet.resource.PathResourceResolver;

import java.io.IOException;

/**
 * Configuration for serving Single Page Application (SPA) static files.
 * Ensures that all client-side routes are properly handled by redirecting to index.html.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Configuration
public class SpaConfig implements WebMvcConfigurer {

    /**
     * Configures resource handlers to serve static content and redirect to index.html
     * for SPA routing.
     */
    @Override
    public void addResourceHandlers(ResourceHandlerRegistry registry) {
        registry.addResourceHandler("/**")
                .addResourceLocations("classpath:/static/")
                .resourceChain(true)
                .addResolver(new PathResourceResolver() {
                    @Override
                    protected Resource getResource(String resourcePath, Resource location) throws IOException {
                        Resource requestedResource = location.createRelative(resourcePath);

                        // If the requested resource exists, serve it (JS, CSS, images, etc.)
                        if (requestedResource.exists() && requestedResource.isReadable()) {
                            return requestedResource;
                        }

                        // For API and backend endpoints, don't redirect to index.html
                        if (resourcePath.startsWith("api/")
                                || resourcePath.startsWith("actuator/")
                                || resourcePath.startsWith("swagger-ui")
                                || resourcePath.startsWith("v3/api-docs")
                                || resourcePath.startsWith("ws/")
                                || resourcePath.startsWith("authenticate")
                                || resourcePath.startsWith("user")
                                || resourcePath.startsWith("container")
                                || resourcePath.startsWith("system")
                                || resourcePath.startsWith("version")) {
                            return null;
                        }

                        // For all other requests (client-side routes), serve index.html
                        return new ClassPathResource("/static/index.html");
                    }
                });
    }
}
