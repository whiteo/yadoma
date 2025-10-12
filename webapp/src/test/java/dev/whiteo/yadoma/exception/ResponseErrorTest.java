package dev.whiteo.yadoma.exception;

import org.junit.jupiter.api.Test;

import java.time.LocalDateTime;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;

class ResponseErrorTest {

    @Test
    void constructor_NoArgs_ShouldCreateEmptyObject() {
        ResponseError responseError = new ResponseError();

        assertNotNull(responseError);
        assertNull(responseError.getTimestamp());
        assertNull(responseError.getStatus());
        assertNull(responseError.getMessage());
        assertNull(responseError.getCode());
    }

    @Test
    void constructor_WithAllArgs_ShouldSetAllFields() {
        LocalDateTime timestamp = LocalDateTime.now();
        Integer status = 404;
        String message = "Not found";
        String code = "ERR_001";

        ResponseError responseError = new ResponseError(timestamp, status, message, code);

        assertEquals(timestamp, responseError.getTimestamp());
        assertEquals(status, responseError.getStatus());
        assertEquals(message, responseError.getMessage());
        assertEquals(code, responseError.getCode());
    }

    @Test
    void builder_ShouldCreateObjectCorrectly() {
        LocalDateTime timestamp = LocalDateTime.now();
        Integer status = 500;
        String message = "Internal server error";
        String code = "ERR_500";

        ResponseError responseError = ResponseError.builder()
                .timestamp(timestamp)
                .status(status)
                .message(message)
                .code(code)
                .build();

        assertEquals(timestamp, responseError.getTimestamp());
        assertEquals(status, responseError.getStatus());
        assertEquals(message, responseError.getMessage());
        assertEquals(code, responseError.getCode());
    }

    @Test
    void settersAndGetters_ShouldWorkCorrectly() {
        ResponseError responseError = new ResponseError();
        LocalDateTime timestamp = LocalDateTime.now();
        Integer status = 400;
        String message = "Bad request";
        String code = "ERR_400";

        responseError.setTimestamp(timestamp);
        responseError.setStatus(status);
        responseError.setMessage(message);
        responseError.setCode(code);

        assertEquals(timestamp, responseError.getTimestamp());
        assertEquals(status, responseError.getStatus());
        assertEquals(message, responseError.getMessage());
        assertEquals(code, responseError.getCode());
    }

    @Test
    void equals_WithSameValues_ShouldReturnTrue() {
        LocalDateTime timestamp = LocalDateTime.now();
        ResponseError error1 = ResponseError.builder()
                .timestamp(timestamp)
                .status(404)
                .message("Not found")
                .code("ERR_404")
                .build();

        ResponseError error2 = ResponseError.builder()
                .timestamp(timestamp)
                .status(404)
                .message("Not found")
                .code("ERR_404")
                .build();

        assertEquals(error1, error2);
    }

    @Test
    void hashCode_WithSameValues_ShouldReturnSameHashCode() {
        LocalDateTime timestamp = LocalDateTime.now();
        ResponseError error1 = ResponseError.builder()
                .timestamp(timestamp)
                .status(404)
                .message("Not found")
                .code("ERR_404")
                .build();

        ResponseError error2 = ResponseError.builder()
                .timestamp(timestamp)
                .status(404)
                .message("Not found")
                .code("ERR_404")
                .build();

        assertEquals(error1.hashCode(), error2.hashCode());
    }
}