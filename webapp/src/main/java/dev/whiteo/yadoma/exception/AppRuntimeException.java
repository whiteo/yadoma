package dev.whiteo.yadoma.exception;

/**
 * Base unchecked exception for application-specific runtime errors.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
public class AppRuntimeException extends RuntimeException {

    /**
     * Constructs a new AppRuntimeException with the specified cause.
     * @param cause the underlying cause of the exception
     */
    public AppRuntimeException(Throwable cause) {
        super(cause);
    }
}