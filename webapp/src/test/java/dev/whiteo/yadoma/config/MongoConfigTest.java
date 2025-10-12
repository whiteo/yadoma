package dev.whiteo.yadoma.config;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;

/**
 * Unit tests for MongoConfig configuration class.
 */
class MongoConfigTest {
    @Test
    void canInstantiate() {
        assertDoesNotThrow(MongoConfig::new);
    }
}