package dev.whiteo.yadoma.exception;

import org.junit.jupiter.api.Test;
import org.springframework.http.HttpStatus;

import static org.junit.jupiter.api.Assertions.assertEquals;

class ExecutionConflictExceptionTest {
    @Test
    void constructor_setsFields() {
        ExecutionConflictException ex = new ExecutionConflictException("conflict");
        assertEquals(ErrorType.CONFLICT, ex.getErrorType());
        assertEquals(HttpStatus.CONFLICT, ex.getStatus());
        assertEquals("conflict", ex.getMessage());
    }
    @Test
    void fullConstructor_setsAllFields() {
        Throwable t = new RuntimeException("err");
        ExecutionConflictException ex = new ExecutionConflictException(ErrorType.FORBIDDEN, HttpStatus.FORBIDDEN, "msg", t);
        assertEquals(ErrorType.FORBIDDEN, ex.getErrorType());
        assertEquals(HttpStatus.FORBIDDEN, ex.getStatus());
        assertEquals("msg", ex.getMessage());
        assertEquals(t, ex.getCause());
    }
}

