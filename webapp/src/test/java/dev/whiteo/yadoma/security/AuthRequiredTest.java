package dev.whiteo.yadoma.security;

import org.junit.jupiter.api.Test;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Test class for AuthRequired annotation.
 */
class AuthRequiredTest {

    @AuthRequired(value = AuthInterceptor.class)
    private static class TestClass {
    }

    @Test
    void annotation_ShouldBePresent() {
        assertTrue(TestClass.class.isAnnotationPresent(AuthRequired.class));
    }

    @Test
    void annotation_ShouldHaveCorrectValue() {
        AuthRequired annotation = TestClass.class.getAnnotation(AuthRequired.class);

        assertEquals(AuthInterceptor.class, annotation.value());
    }

    @Test
    void annotation_ShouldHaveCorrectRetentionPolicy() {
        Retention retention = AuthRequired.class.getAnnotation(Retention.class);

        assertEquals(RetentionPolicy.RUNTIME, retention.value());
    }

    @Test
    void annotation_ShouldHaveCorrectTarget() {
        Target target = AuthRequired.class.getAnnotation(Target.class);

        assertEquals(1, target.value().length);
        assertEquals(ElementType.TYPE, target.value()[0]);
    }
}
