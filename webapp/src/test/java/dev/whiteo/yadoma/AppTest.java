package dev.whiteo.yadoma;

import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;

import java.util.TimeZone;

import static org.junit.jupiter.api.Assertions.assertEquals;

/**
 * Test class for the main App class.
 */
@SpringBootTest
class AppTest {

    @Test
    void contextLoads() {
        /*
         * This test method is intentionally empty.
         * It serves as a basic smoke test to verify that the Spring Boot application context
         * can be loaded successfully without any configuration errors or bean creation failures.
         * The @SpringBootTest annotation automatically triggers the application context loading,
         * and if the context fails to load, this test will fail.
         * No additional assertions are needed as the framework handles the validation.
         */
    }

    @Test
    void mainMethod_SetsTimezoneCorrectly() {
        App.main(new String[]{});
        assertEquals("Europe/Berlin", TimeZone.getDefault().getID());
    }
}
