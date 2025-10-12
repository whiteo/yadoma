package dev.whiteo.yadoma.exception;

/**
 * Enumeration of error types used for error handling and response codes.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
public enum ErrorType {
    UNDEFINED_ERROR,
    BAD_REQUEST,
    NOT_FOUND,
    AUTHORIZATION_ERROR,
    REQUEST_MESSAGE_CONVERSION_ERROR,
    REQUEST_INVALID_MESSAGE,
    REQUEST_TYPE_MISMATCH,
    REQUEST_MISSING_PARAMETER,
    UNSUPPORTED_OPERATION,
    FORBIDDEN,
    CONFLICT,
    REQUEST_ROUTE_ERROR
}