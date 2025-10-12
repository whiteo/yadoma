package dev.whiteo.yadoma.exception;

import org.junit.jupiter.api.Test;
import org.springframework.http.converter.HttpMessageConversionException;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

class ExceptionHelperTest {

    @Test
    void getRootCauseMessage_returnsRootCauseMessage() {
        RuntimeException root = new RuntimeException("root-cause");
        Exception ex = new Exception("wrapper", root);

        String msg = ExceptionHelper.getRootCauseMessage(ex);
        assertTrue(msg.contains("root-cause"));
    }

    @Test
    void getJsonExceptionMessage_null_returnsDefault() {
        assertEquals(ExceptionHelper.UNDEFINED_FORMAT_ERROR, ExceptionHelper.getJsonExceptionMessage(null));
    }

    @Test
    void getJsonExceptionMessage_simpleException_returnsRootMessage() {
        HttpMessageConversionException ex = new HttpMessageConversionException("simple error");
        String msg = ExceptionHelper.getJsonExceptionMessage(ex);
        assertTrue(msg.contains("simple error"));
    }
}
