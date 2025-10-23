package dev.whiteo.yadoma.config;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;

class MongoConfigTest {
    @Test
    void canInstantiate() {
        assertDoesNotThrow(MongoConfig::new);
    }
}