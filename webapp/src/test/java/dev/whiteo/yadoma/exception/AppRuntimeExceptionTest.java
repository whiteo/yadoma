package dev.whiteo.yadoma.exception;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Test class for AppRuntimeException.
 */
class AppRuntimeExceptionTest {

    @Test
    void constructor_ShouldCreateExceptionWithCause() {
        RuntimeException cause = new RuntimeException("Original cause");

        AppRuntimeException exception = new AppRuntimeException(cause);

        assertEquals(cause, exception.getCause());
        assertEquals("Original cause", exception.getCause().getMessage());
    }

    @Test
    void constructor_ShouldHandleNullCause() {
        AppRuntimeException exception = new AppRuntimeException(null);

        assertNull(exception.getCause());
    }

    @Test
    void constructor_ShouldCreateExceptionWithNestedCause() {
        Exception rootCause = new Exception("Root cause");
        RuntimeException intermediateCause = new RuntimeException("Intermediate cause", rootCause);

        AppRuntimeException exception = new AppRuntimeException(intermediateCause);

        assertEquals(intermediateCause, exception.getCause());
        assertEquals("Intermediate cause", exception.getCause().getMessage());
        assertEquals(rootCause, exception.getCause().getCause());
    }

    @Test
    void inheritance_ShouldExtendRuntimeException() {
        AppRuntimeException exception = new AppRuntimeException(new RuntimeException());

        assertTrue(exception instanceof RuntimeException);
    }
}
