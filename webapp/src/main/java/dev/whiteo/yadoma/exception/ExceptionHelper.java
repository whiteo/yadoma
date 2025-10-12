package dev.whiteo.yadoma.exception;

import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.util.ClassUtil;
import lombok.NoArgsConstructor;
import org.springframework.http.converter.HttpMessageConversionException;

import java.util.stream.Collectors;

/**
 * Utility class for extracting and formatting exception messages.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@NoArgsConstructor
public class ExceptionHelper {

    public static final String UNDEFINED_FORMAT_ERROR = "Undefined format error";

    /**
     * Returns the root cause message of the given exception.
     * @param ex the exception
     * @return root cause message string
     */
    public static String getRootCauseMessage(Throwable ex) {
        Throwable rootCause = ClassUtil.getRootCause(ex);
        return ClassUtil.exceptionMessage(rootCause);
    }

    /**
     * Returns a formatted message for JSON conversion exceptions.
     * @param ex the HttpMessageConversionException
     * @return formatted error message string
     */
    public static String getJsonExceptionMessage(HttpMessageConversionException ex) {
        if (ex == null) {
            return UNDEFINED_FORMAT_ERROR;
        }

        String message = getRootCauseMessage(ex);
        Throwable cause = ex.getCause();
        if (cause instanceof JsonMappingException error) {
            String field = error.getPath()
                    .stream()
                    .map(JsonMappingException.Reference::getFieldName)
                    .collect(Collectors.joining("."));
            return "Format error: " + message + " link to '" + field + "'";
        }
        return message;
    }
}