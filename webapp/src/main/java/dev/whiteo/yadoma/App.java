package dev.whiteo.yadoma;

import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import java.util.TimeZone;

/**
 * Main entry point for the Yadoma Spring Boot application.
 * Sets the default timezone and starts the application context.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Slf4j
@SpringBootApplication
public class App {

    /**
     * Starts the Spring Boot application.
     * @param args command-line arguments
     */
    public static void main(String[] args) {
        TimeZone.setDefault(TimeZone.getTimeZone("Europe/Berlin"));
        SpringApplication.run(App.class, args);
    }
}