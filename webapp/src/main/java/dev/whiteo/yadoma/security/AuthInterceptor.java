package dev.whiteo.yadoma.security;

import io.micrometer.common.util.StringUtils;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.servlet.ModelAndView;

/**
 * Interceptor for authenticating requests using JWT tokens.
 * Extracts user information from token and sets it in ThreadLocal.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Component
@RequiredArgsConstructor
@Order(Ordered.HIGHEST_PRECEDENCE)
public class AuthInterceptor implements HandlerInterceptor {

    private final TokenInteract tokenInteract;
    private final ThreadLocal<String> users = new ThreadLocal<>();

    /**
     * Checks JWT token validity and sets user in ThreadLocal before handling request.
     */
    @Override
    public boolean preHandle(HttpServletRequest request, @NonNull HttpServletResponse response,
                             @NonNull Object handler) throws Exception {
        String path = request.getRequestURI();
        if (path.endsWith("/version")
                || path.endsWith("/api/v1/authenticate")
                || path.endsWith("/api/v1/user/create")
                || path.contains("/v3/api-docs")
                || path.contains("/swagger-ui/")) {
            return true;
        }

        String token = tokenInteract.getToken(request);
        if (tokenInteract.validateToken(token)) {
            String user = tokenInteract.getUserId(token);
            if (StringUtils.isNotBlank(user)) {
                users.set(user);
                return true;
            }
        }
        response.sendError(HttpServletResponse.SC_UNAUTHORIZED, "Unauthorized");
        return false;
    }

    /**
     * Returns the userId extracted from the JWT token for the current request thread.
     */
    public String getUserId() {
        return users.get();
    }

    /**
     * Removes user information from ThreadLocal after request is handled.
     */
    @Override
    public void postHandle(@NonNull HttpServletRequest request, @NonNull HttpServletResponse response,
                           @NonNull Object handler, ModelAndView modelAndView) {
        users.remove();
    }
}