package dev.whiteo.yadoma.exception;

import lombok.Getter;
import org.springframework.http.HttpStatus;

/**
 * Exception thrown when an entity is not found in the database.
 * Contains error type and HTTP status.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Getter
public class NotFoundObjectException extends RuntimeException {

    private final ErrorType errorType;
    private final HttpStatus status;

    /**
     * Constructs a new NotFoundObjectException with default error type and status.
     * @param message error message
     */
    public NotFoundObjectException(String message) {
        this(ErrorType.NOT_FOUND, HttpStatus.NOT_FOUND, message, null);
    }

    /**
     * Constructs a new NotFoundObjectException with custom error type, status, message, and cause.
     * @param errorType type of error
     * @param status HTTP status
     * @param message error message
     * @param error underlying cause
     */
    public NotFoundObjectException(ErrorType errorType, HttpStatus status, String message, Throwable error) {
        super(message, error);
        this.errorType = errorType;
        this.status = status;
    }
}