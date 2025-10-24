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
}
