package dev.whiteo.yadoma;

import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;

import java.util.TimeZone;

import static org.junit.jupiter.api.Assertions.assertEquals;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.MOCK)
class AppTest {

    @Test
    void contextLoads() {
    }

    @Test
    void mainMethod_SetsTimezoneCorrectly() {
        TimeZone.setDefault(TimeZone.getTimeZone("UTC"));
        assertEquals("UTC", TimeZone.getDefault().getID());
        TimeZone.setDefault(TimeZone.getTimeZone("Europe/Berlin"));
        assertEquals("Europe/Berlin", TimeZone.getDefault().getID());
    }
}
