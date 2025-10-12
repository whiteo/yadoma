package dev.whiteo.yadoma.exception;

import org.junit.jupiter.api.Test;
import org.springframework.http.HttpStatus;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

class NotFoundObjectExceptionTest {

    @Test
    void constructor_WithMessage_ShouldSetMessage() {
        String message = "Object not found";

        NotFoundObjectException exception = new NotFoundObjectException(message);

        assertEquals(message, exception.getMessage());
        assertEquals(ErrorType.NOT_FOUND, exception.getErrorType());
        assertEquals(HttpStatus.NOT_FOUND, exception.getStatus());
    }

    @Test
    void constructor_WithAllParameters_ShouldSetAllFields() {
        String message = "Object not found";
        Throwable cause = new RuntimeException("Root cause");
        ErrorType errorType = ErrorType.NOT_FOUND;
        HttpStatus status = HttpStatus.NOT_FOUND;

        NotFoundObjectException exception = new NotFoundObjectException(errorType, status, message, cause);

        assertEquals(message, exception.getMessage());
        assertEquals(cause, exception.getCause());
        assertEquals(errorType, exception.getErrorType());
        assertEquals(status, exception.getStatus());
    }

    @Test
    void isRuntimeException_ShouldExtendRuntimeException() {
        NotFoundObjectException exception = new NotFoundObjectException("Test");

        assertTrue(exception instanceof RuntimeException);
    }

    @Test
    void getErrorType_ShouldReturnCorrectType() {
        NotFoundObjectException exception = new NotFoundObjectException("Test message");

        assertEquals(ErrorType.NOT_FOUND, exception.getErrorType());
    }

    @Test
    void getStatus_ShouldReturnCorrectStatus() {
        NotFoundObjectException exception = new NotFoundObjectException("Test message");

        assertEquals(HttpStatus.NOT_FOUND, exception.getStatus());
    }
}