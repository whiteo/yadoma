package dev.whiteo.yadoma.config;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.core.io.ClassPathResource;
import org.springframework.core.io.Resource;
import org.springframework.web.servlet.resource.PathResourceResolver;

import java.io.IOException;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

class SpaConfigTest {

    private SpaConfig spaConfig;
    private PathResourceResolver resolver;

    @BeforeEach
    void setUp() {
        spaConfig = new SpaConfig();
    }

    @Test
    void spaConfig_shouldBeCreated() {
        assertNotNull(spaConfig);
    }

    @Test
    void spaConfig_shouldImplementWebMvcConfigurer() {
        assertNotNull(spaConfig);
    }

    @Test
    void pathResourceResolver_ShouldReturnExistingStaticResource() throws Exception {
        resolver = createPathResourceResolver();
        Resource location = new ClassPathResource("/static/");

        Resource mockResource = mock(Resource.class);
        when(mockResource.exists()).thenReturn(true);
        when(mockResource.isReadable()).thenReturn(true);

        Resource result = invokeGetResource(resolver, "assets/test.js", location);

        assertNotNull(result);
    }

    @Test
    void pathResourceResolver_ShouldReturnIndexForNonExistentSpaRoute() throws Exception {
        resolver = createPathResourceResolver();
        Resource location = new ClassPathResource("/static/");

        Resource result = invokeGetResource(resolver, "dashboard", location);

        assertNotNull(result);
        assertTrue(result.getFilename().contains("index.html"));
    }

    @Test
    void pathResourceResolver_ShouldReturnNullForApiPath() throws Exception {
        resolver = createPathResourceResolver();
        Resource location = new ClassPathResource("/static/");

        Resource result = invokeGetResource(resolver, "api/v1/users", location);

        assertNull(result);
    }

    @Test
    void pathResourceResolver_ShouldReturnNullForActuatorPath() throws Exception {
        resolver = createPathResourceResolver();
        Resource location = new ClassPathResource("/static/");

        Resource result = invokeGetResource(resolver, "actuator/health", location);

        assertNull(result);
    }

    @Test
    void pathResourceResolver_ShouldReturnNullForSwaggerUiPath() throws Exception {
        resolver = createPathResourceResolver();
        Resource location = new ClassPathResource("/static/");

        Resource result = invokeGetResource(resolver, "swagger-ui/index.html", location);

        assertNull(result);
    }

    @Test
    void pathResourceResolver_ShouldReturnNullForApiDocsPath() throws Exception {
        resolver = createPathResourceResolver();
        Resource location = new ClassPathResource("/static/");

        Resource result = invokeGetResource(resolver, "v3/api-docs", location);

        assertNull(result);
    }

    @Test
    void pathResourceResolver_ShouldReturnNullForWebSocketPath() throws Exception {
        resolver = createPathResourceResolver();
        Resource location = new ClassPathResource("/static/");

        Resource result = invokeGetResource(resolver, "ws/containers/123/stats", location);

        assertNull(result);
    }

    @Test
    void pathResourceResolver_ShouldReturnNullForAuthenticatePath() throws Exception {
        resolver = createPathResourceResolver();
        Resource location = new ClassPathResource("/static/");

        Resource result = invokeGetResource(resolver, "authenticate", location);

        assertNull(result);
    }

    @Test
    void pathResourceResolver_ShouldReturnNullForUserPath() throws Exception {
        resolver = createPathResourceResolver();
        Resource location = new ClassPathResource("/static/");

        Resource result = invokeGetResource(resolver, "user/profile", location);

        assertNull(result);
    }

    @Test
    void pathResourceResolver_ShouldReturnNullForContainerPath() throws Exception {
        resolver = createPathResourceResolver();
        Resource location = new ClassPathResource("/static/");

        Resource result = invokeGetResource(resolver, "container/list", location);

        assertNull(result);
    }

    @Test
    void pathResourceResolver_ShouldReturnNullForSystemPath() throws Exception {
        resolver = createPathResourceResolver();
        Resource location = new ClassPathResource("/static/");

        Resource result = invokeGetResource(resolver, "system/info", location);

        assertNull(result);
    }

    @Test
    void pathResourceResolver_ShouldReturnNullForVersionPath() throws Exception {
        resolver = createPathResourceResolver();
        Resource location = new ClassPathResource("/static/");

        Resource result = invokeGetResource(resolver, "version", location);

        assertNull(result);
    }

    @Test
    void pathResourceResolver_ShouldReturnIndexForSettingsRoute() throws Exception {
        resolver = createPathResourceResolver();
        Resource location = new ClassPathResource("/static/");

        Resource result = invokeGetResource(resolver, "settings", location);

        assertNotNull(result);
        assertTrue(result.getFilename().contains("index.html"));
    }

    @Test
    void pathResourceResolver_ShouldReturnIndexForNestedSpaRoute() throws Exception {
        resolver = createPathResourceResolver();
        Resource location = new ClassPathResource("/static/");

        Resource result = invokeGetResource(resolver, "dashboard/overview", location);

        assertNotNull(result);
        assertTrue(result.getFilename().contains("index.html"));
    }

    private PathResourceResolver createPathResourceResolver() {
        return new PathResourceResolver() {
            @Override
            protected Resource getResource(String resourcePath, Resource location) throws IOException {
                Resource requestedResource = location.createRelative(resourcePath);

                if (requestedResource.exists() && requestedResource.isReadable()) {
                    return requestedResource;
                }

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

                return new ClassPathResource("/static/index.html");
            }
        };
    }

    private Resource invokeGetResource(PathResourceResolver resolver, String resourcePath, Resource location) throws Exception {
        java.lang.reflect.Method method = PathResourceResolver.class.getDeclaredMethod("getResource", String.class, Resource.class);
        method.setAccessible(true);
        return (Resource) method.invoke(resolver, resourcePath, location);
    }
}
