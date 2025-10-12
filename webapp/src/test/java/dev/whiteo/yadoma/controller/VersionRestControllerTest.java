package dev.whiteo.yadoma.controller;

import org.junit.jupiter.api.Test;
import org.springframework.test.util.ReflectionTestUtils;

import static org.junit.jupiter.api.Assertions.assertEquals;

class VersionRestControllerTest {

    @Test
    void getAppVersion_shouldReturnInjectedValue() {
        VersionRestController controller = new VersionRestController();
        ReflectionTestUtils.setField(controller, "appVersion", "1.2.3");
        String v = controller.getAppVersion();
        assertEquals("1.2.3", v);
    }
}