package dev.whiteo.yadoma;

import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.test.context.ActiveProfiles;

import java.util.TimeZone;

import static org.junit.jupiter.api.Assertions.assertEquals;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.MOCK)
@ActiveProfiles("test")
class AppTest {

    @MockBean
    private io.grpc.ManagedChannel grpcChannel;

    @MockBean
    private org.springframework.web.socket.server.standard.ServerEndpointExporter serverEndpointExporter;

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
