package dev.whiteo.yadoma.exception;

import jakarta.validation.ConstraintViolationException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.TypeMismatchException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.http.converter.HttpMessageConversionException;
import org.springframework.http.converter.HttpMessageNotReadableException;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.validation.BindingResult;
import org.springframework.validation.FieldError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.ServletRequestBindingException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestControllerAdvice;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.method.annotation.MethodArgumentTypeMismatchException;

import java.time.LocalDateTime;
import java.util.List;

import static dev.whiteo.yadoma.exception.ExceptionHelper.getJsonExceptionMessage;
import static dev.whiteo.yadoma.exception.ExceptionHelper.getRootCauseMessage;
import static java.util.Optional.ofNullable;
import static org.springframework.http.HttpStatus.BAD_REQUEST;
import static org.springframework.http.HttpStatus.CONFLICT;
import static org.springframework.http.HttpStatus.INTERNAL_SERVER_ERROR;
import static org.springframework.http.HttpStatus.NOT_FOUND;
import static org.springframework.http.HttpStatus.UNAUTHORIZED;

/**
 * Global exception handler for REST API.
 * Handles various exceptions and returns standardized error responses.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Slf4j
@RestControllerAdvice
@RequiredArgsConstructor
public class RestExceptionHandler {

    /**
     * Builds a standardized error response for REST API.
     * @param status HTTP status
     * @param code error type
     * @param msg error message
     * @param ex exception
     * @return ResponseEntity with ResponseError
     */
    public static ResponseEntity<ResponseError> buildResponse(HttpStatus status, ErrorType code, String msg,
                                                              Throwable ex) {
        if (NOT_FOUND.equals(status) || CONFLICT.equals(status)) {
            log.warn(msg);
        } else {
            log.error(msg, ex);
        }
        return ResponseEntity.status(status.value())
                .body(ResponseError.builder()
                        .code(code.name())
                        .message(msg)
                        .timestamp(LocalDateTime.now())
                        .status(status.value())
                        .build());
    }

    /**
     * Handles all uncaught exceptions and returns an undefined error response.
     */
    @ExceptionHandler(Throwable.class)
    @ResponseStatus(INTERNAL_SERVER_ERROR)
    public ResponseEntity<ResponseError> handleUndefinedException(Throwable ex) {
        String message = "Undefined error: " + getRootCauseMessage(ex);
        return buildResponse(INTERNAL_SERVER_ERROR, ErrorType.UNDEFINED_ERROR, message, ex);
    }

    /**
     * Handles HttpClientErrorException (authorization errors).
     */
    @ResponseStatus(UNAUTHORIZED)
    @ExceptionHandler(HttpClientErrorException.class)
    protected ResponseEntity<ResponseError> handleHttpClientErrorException(HttpClientErrorException ex) {
        return buildResponse(UNAUTHORIZED, ErrorType.AUTHORIZATION_ERROR, ex.getMessage(), ex);
    }

    /**
     * Handles BadCredentialsException (authorization errors).
     */
    @ResponseStatus(UNAUTHORIZED)
    @ExceptionHandler(BadCredentialsException.class)
    protected ResponseEntity<ResponseError> handleBadCredentialsException(BadCredentialsException ex) {
        return buildResponse(UNAUTHORIZED, ErrorType.AUTHORIZATION_ERROR, ex.getMessage(), ex);
    }

    /**
     * Handles ExecutionConflictException (conflict errors).
     */
    @ResponseStatus(CONFLICT)
    @ExceptionHandler(ExecutionConflictException.class)
    protected ResponseEntity<ResponseError> handleExecutionConflictException(ExecutionConflictException ex) {
        return buildResponse(CONFLICT, ErrorType.CONFLICT, ex.getMessage(), ex);
    }

    /**
     * Handles NotFoundObjectException (not found errors).
     */
    @ResponseStatus(NOT_FOUND)
    @ExceptionHandler(NotFoundObjectException.class)
    protected ResponseEntity<ResponseError> handleNotFoundObjectException(NotFoundObjectException ex) {
        return buildResponse(NOT_FOUND, ErrorType.NOT_FOUND, ex.getMessage(), ex);
    }

    /**
     * Handles MethodArgumentTypeMismatchException (bad request errors).
     */
    @ResponseStatus(BAD_REQUEST)
    @ExceptionHandler(MethodArgumentTypeMismatchException.class)
    protected ResponseEntity<ResponseError> handleMethodArgumentTypeMismatchException(
            MethodArgumentTypeMismatchException ex) {
        return buildResponse(BAD_REQUEST, ErrorType.BAD_REQUEST, ex.getMessage(), ex);
    }

    /**
     * Handles HttpMessageNotReadableException (bad request errors).
     */
    @ResponseStatus(BAD_REQUEST)
    @ExceptionHandler(HttpMessageNotReadableException.class)
    public ResponseEntity<ResponseError> handleHttpMessageNotReadableException(HttpMessageNotReadableException ex) {
        return buildResponse(BAD_REQUEST, ErrorType.BAD_REQUEST, ex.getMessage(), ex);
    }

    /**
     * Handles HttpMessageConversionException (message conversion errors).
     */
    @ResponseStatus(BAD_REQUEST)
    @ExceptionHandler(HttpMessageConversionException.class)
    public ResponseEntity<ResponseError> handleMessageConversionException(HttpMessageConversionException ex) {
        return buildResponse(BAD_REQUEST, ErrorType.REQUEST_MESSAGE_CONVERSION_ERROR,
                getJsonExceptionMessage(ex), ex);
    }

    /**
     * Handles ServletRequestBindingException (missing parameter errors).
     */
    @ResponseStatus(BAD_REQUEST)
    @ExceptionHandler(ServletRequestBindingException.class)
    public ResponseEntity<ResponseError> handleServletRequestBinding(ServletRequestBindingException ex) {
        return buildResponse(BAD_REQUEST, ErrorType.REQUEST_MISSING_PARAMETER,
                getRootCauseMessage(ex), ex);
    }

    /**
     * Handles ConstraintViolationException (missing parameter errors).
     */
    @ResponseStatus(BAD_REQUEST)
    @ExceptionHandler(ConstraintViolationException.class)
    public ResponseEntity<ResponseError> handleConstraintViolation(ConstraintViolationException ex) {
        return buildResponse(BAD_REQUEST, ErrorType.REQUEST_MISSING_PARAMETER,
                getRootCauseMessage(ex), ex);
    }

    /**
     * Handles TypeMismatchException (type mismatch errors).
     */
    @ResponseStatus(BAD_REQUEST)
    @ExceptionHandler(TypeMismatchException.class)
    public ResponseEntity<ResponseError> handleTypeMismatch(TypeMismatchException ex) {
        String requiredType = ofNullable(ex.getRequiredType()).map(Class::getSimpleName).orElse("<empty>");
        String message = getRootCauseMessage(ex) + ", required type '" + requiredType + "'";
        return buildResponse(BAD_REQUEST, ErrorType.REQUEST_TYPE_MISMATCH, message, ex);
    }

    /**
     * Handles UnsupportedOperationException (unsupported operation errors).
     */
    @ResponseStatus(BAD_REQUEST)
    @ExceptionHandler(UnsupportedOperationException.class)
    public ResponseEntity<ResponseError> handleUnsupportedOperation(UnsupportedOperationException ex) {
        return buildResponse(BAD_REQUEST, ErrorType.UNSUPPORTED_OPERATION, getRootCauseMessage(ex), ex);
    }

    /**
     * Handles MethodArgumentNotValidException (validation errors).
     */
    @ResponseBody
    @ResponseStatus(BAD_REQUEST)
    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<ResponseError> validationError(MethodArgumentNotValidException ex) {
        BindingResult result = ex.getBindingResult();
        List<FieldError> fieldErrors = result.getFieldErrors();
        String errors = fieldErrors.stream()
                .map(error -> error.getField() + ": " + error.getDefaultMessage()).findFirst().orElse("ops!");
        return buildResponse(BAD_REQUEST, ErrorType.BAD_REQUEST, errors, ex);
    }
}