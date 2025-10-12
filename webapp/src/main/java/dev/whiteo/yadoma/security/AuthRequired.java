package dev.whiteo.yadoma.security;

import org.springframework.web.servlet.HandlerInterceptor;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

/**
 * Annotation to mark controllers or classes that require authentication.
 * Specifies the HandlerInterceptor to use for authorization.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Target({ElementType.TYPE})
@Retention(RetentionPolicy.RUNTIME)
public @interface AuthRequired {

    /**
     * The interceptor class to use for authentication.
     * @return HandlerInterceptor class
     */
    Class<? extends HandlerInterceptor> value();
}