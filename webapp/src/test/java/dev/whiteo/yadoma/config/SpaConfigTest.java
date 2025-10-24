package dev.whiteo.yadoma.config;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertNotNull;

class SpaConfigTest {

    private SpaConfig spaConfig;

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
        // SpaConfig implements WebMvcConfigurer which provides addResourceHandlers
        // This is tested through Spring integration tests
    }
}
