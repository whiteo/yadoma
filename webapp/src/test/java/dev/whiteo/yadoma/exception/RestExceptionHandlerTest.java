package dev.whiteo.yadoma.exception;

import jakarta.validation.ConstraintViolationException;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.junit.jupiter.MockitoExtension;
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
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.method.annotation.MethodArgumentTypeMismatchException;

import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class RestExceptionHandlerTest {

    @InjectMocks
    private RestExceptionHandler restExceptionHandler;

    @Test
    void buildResponse_ShouldCreateCorrectResponseEntity() {
        HttpStatus status = HttpStatus.BAD_REQUEST;
        ErrorType errorType = ErrorType.BAD_REQUEST;
        String message = "Test error message";
        Exception exception = new RuntimeException("Test exception");

        ResponseEntity<ResponseError> response = RestExceptionHandler.buildResponse(status, errorType, message, exception);

        assertNotNull(response);
        assertEquals(HttpStatus.BAD_REQUEST, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(errorType.name(), response.getBody().getCode());
        assertEquals(message, response.getBody().getMessage());
        assertEquals(400, response.getBody().getStatus());
        assertNotNull(response.getBody().getTimestamp());
    }

    @Test
    void handleUndefinedException_ShouldReturnInternalServerError() {
        RuntimeException exception = new RuntimeException("Unexpected error");

        ResponseEntity<ResponseError> response = restExceptionHandler.handleUndefinedException(exception);

        assertNotNull(response);
        assertEquals(HttpStatus.INTERNAL_SERVER_ERROR, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(ErrorType.UNDEFINED_ERROR.name(), response.getBody().getCode());
        assertTrue(response.getBody().getMessage().contains("Undefined error"));
        assertEquals(500, response.getBody().getStatus());
    }

    @Test
    void handleHttpClientErrorException_ShouldReturnUnauthorized() {
        HttpClientErrorException exception = new HttpClientErrorException(HttpStatus.UNAUTHORIZED, "Unauthorized access");

        ResponseEntity<ResponseError> response = restExceptionHandler.handleHttpClientErrorException(exception);

        assertNotNull(response);
        assertEquals(HttpStatus.UNAUTHORIZED, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(ErrorType.AUTHORIZATION_ERROR.name(), response.getBody().getCode());
        assertTrue(response.getBody().getMessage().contains("Unauthorized"));
        assertEquals(401, response.getBody().getStatus());
    }

    @Test
    void handleBadCredentialsException_ShouldReturnUnauthorized() {
        BadCredentialsException exception = new BadCredentialsException("Bad credentials");

        ResponseEntity<ResponseError> response = restExceptionHandler.handleBadCredentialsException(exception);

        assertNotNull(response);
        assertEquals(HttpStatus.UNAUTHORIZED, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(ErrorType.AUTHORIZATION_ERROR.name(), response.getBody().getCode());
        assertEquals("Bad credentials", response.getBody().getMessage());
        assertEquals(401, response.getBody().getStatus());
    }

    @Test
    void handleExecutionConflictException_ShouldReturnConflict() {
        ExecutionConflictException exception = new ExecutionConflictException("Conflict occurred");

        ResponseEntity<ResponseError> response = restExceptionHandler.handleExecutionConflictException(exception);

        assertNotNull(response);
        assertEquals(HttpStatus.CONFLICT, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(ErrorType.CONFLICT.name(), response.getBody().getCode());
        assertEquals("Conflict occurred", response.getBody().getMessage());
        assertEquals(409, response.getBody().getStatus());
    }

    @Test
    void handleNotFoundObjectException_ShouldReturnNotFound() {
        NotFoundObjectException exception = new NotFoundObjectException("Object not found");

        ResponseEntity<ResponseError> response = restExceptionHandler.handleNotFoundObjectException(exception);

        assertNotNull(response);
        assertEquals(HttpStatus.NOT_FOUND, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(ErrorType.NOT_FOUND.name(), response.getBody().getCode());
        assertEquals("Object not found", response.getBody().getMessage());
        assertEquals(404, response.getBody().getStatus());
    }

    @Test
    void handleMethodArgumentTypeMismatchException_ShouldReturnBadRequest() {
        MethodArgumentTypeMismatchException exception = mock(MethodArgumentTypeMismatchException.class);
        when(exception.getMessage()).thenReturn("Type mismatch error");

        ResponseEntity<ResponseError> response = restExceptionHandler.handleMethodArgumentTypeMismatchException(exception);

        assertNotNull(response);
        assertEquals(HttpStatus.BAD_REQUEST, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(ErrorType.BAD_REQUEST.name(), response.getBody().getCode());
        assertEquals("Type mismatch error", response.getBody().getMessage());
        assertEquals(400, response.getBody().getStatus());
    }

    @Test
    void handleHttpMessageNotReadableException_ShouldReturnBadRequest() {
        HttpMessageNotReadableException exception = mock(HttpMessageNotReadableException.class);
        when(exception.getMessage()).thenReturn("Message not readable");

        ResponseEntity<ResponseError> response = restExceptionHandler.handleHttpMessageNotReadableException(exception);

        assertNotNull(response);
        assertEquals(HttpStatus.BAD_REQUEST, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(ErrorType.BAD_REQUEST.name(), response.getBody().getCode());
        assertEquals("Message not readable", response.getBody().getMessage());
        assertEquals(400, response.getBody().getStatus());
    }

    @Test
    void handleMessageConversionException_ShouldReturnBadRequest() {
        HttpMessageConversionException exception = mock(HttpMessageConversionException.class);
        when(exception.getMessage()).thenReturn("Conversion error");

        ResponseEntity<ResponseError> response = restExceptionHandler.handleMessageConversionException(exception);

        assertNotNull(response);
        assertEquals(HttpStatus.BAD_REQUEST, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(ErrorType.REQUEST_MESSAGE_CONVERSION_ERROR.name(), response.getBody().getCode());
        assertEquals(400, response.getBody().getStatus());
    }

    @Test
    void handleServletRequestBinding_ShouldReturnBadRequest() {
        ServletRequestBindingException exception = mock(ServletRequestBindingException.class);
        when(exception.getMessage()).thenReturn("Missing parameter");

        ResponseEntity<ResponseError> response = restExceptionHandler.handleServletRequestBinding(exception);

        assertNotNull(response);
        assertEquals(HttpStatus.BAD_REQUEST, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(ErrorType.REQUEST_MISSING_PARAMETER.name(), response.getBody().getCode());
        assertEquals(400, response.getBody().getStatus());
    }

    @Test
    void handleConstraintViolation_ShouldReturnBadRequest() {
        ConstraintViolationException exception = mock(ConstraintViolationException.class);
        when(exception.getMessage()).thenReturn("Constraint violation");

        ResponseEntity<ResponseError> response = restExceptionHandler.handleConstraintViolation(exception);

        assertNotNull(response);
        assertEquals(HttpStatus.BAD_REQUEST, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(ErrorType.REQUEST_MISSING_PARAMETER.name(), response.getBody().getCode());
        assertEquals(400, response.getBody().getStatus());
    }

    @Test
    void handleTypeMismatch_ShouldReturnBadRequest() {
        TypeMismatchException exception = mock(TypeMismatchException.class);
        when(exception.getMessage()).thenReturn("Type mismatch");
        when(exception.getRequiredType()).thenReturn((Class) String.class);

        ResponseEntity<ResponseError> response = restExceptionHandler.handleTypeMismatch(exception);

        assertNotNull(response);
        assertEquals(HttpStatus.BAD_REQUEST, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(ErrorType.REQUEST_TYPE_MISMATCH.name(), response.getBody().getCode());
        assertTrue(response.getBody().getMessage().contains("String"));
        assertEquals(400, response.getBody().getStatus());
    }

    @Test
    void handleUnsupportedOperation_ShouldReturnBadRequest() {
        UnsupportedOperationException exception = new UnsupportedOperationException("Operation not supported");

        ResponseEntity<ResponseError> response = restExceptionHandler.handleUnsupportedOperation(exception);

        assertNotNull(response);
        assertEquals(HttpStatus.BAD_REQUEST, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(ErrorType.UNSUPPORTED_OPERATION.name(), response.getBody().getCode());
        assertEquals("Operation not supported", response.getBody().getMessage());
        assertEquals(400, response.getBody().getStatus());
    }

    @Test
    void validationError_ShouldReturnBadRequest() {
        MethodArgumentNotValidException exception = mock(MethodArgumentNotValidException.class);
        BindingResult bindingResult = mock(BindingResult.class);
        FieldError fieldError = new FieldError("testObject", "testField", "Test validation error");

        when(exception.getBindingResult()).thenReturn(bindingResult);
        when(bindingResult.getFieldErrors()).thenReturn(List.of(fieldError));

        ResponseEntity<ResponseError> response = restExceptionHandler.validationError(exception);

        assertNotNull(response);
        assertEquals(HttpStatus.BAD_REQUEST, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals(ErrorType.BAD_REQUEST.name(), response.getBody().getCode());
        assertTrue(response.getBody().getMessage().contains("testField"));
        assertTrue(response.getBody().getMessage().contains("Test validation error"));
        assertEquals(400, response.getBody().getStatus());
    }
}