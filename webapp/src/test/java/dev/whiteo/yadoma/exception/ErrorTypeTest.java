package dev.whiteo.yadoma.exception;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.junit.jupiter.api.Assertions.fail;

class ErrorTypeTest {

    @Test
    void enumContainsExpectedValues() {
        ErrorType[] values = ErrorType.values();
        assertTrue(values.length > 0);
        assertTrue(java.util.EnumSet.allOf(ErrorType.class).contains(ErrorType.NOT_FOUND));
    }

    @Test
    void allEnumValuesAreUsableInSwitch() {
        for (ErrorType type : ErrorType.values()) {
            switch (type) {
                case UNDEFINED_ERROR, BAD_REQUEST, NOT_FOUND, AUTHORIZATION_ERROR,
                     REQUEST_MESSAGE_CONVERSION_ERROR, REQUEST_INVALID_MESSAGE,
                     REQUEST_TYPE_MISMATCH, REQUEST_MISSING_PARAMETER,
                     UNSUPPORTED_OPERATION, FORBIDDEN, CONFLICT, REQUEST_ROUTE_ERROR ->
                    assertNotNull(type);
                default -> fail("Unknown ErrorType: " + type);
            }
        }
    }
}
