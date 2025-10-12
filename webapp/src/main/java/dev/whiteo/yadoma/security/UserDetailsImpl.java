package dev.whiteo.yadoma.security;

import lombok.Builder;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;

import java.util.Collection;
import java.util.Collections;

/**
 * Implementation of UserDetails for Spring Security authentication.
 * Stores user ID, name, and password hash.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Builder
public class UserDetailsImpl implements UserDetails {

    private String id;

    /**
     * Returns authorities granted to the user (empty for this implementation).
     */
    @Override
    public Collection<? extends GrantedAuthority> getAuthorities() {
        return Collections.emptyList();
    }

    /**
     * Returns the user's password (empty for this implementation).
     */
    @Override
    public String getPassword() {
        return "";
    }

    /**
     * Returns the user's username (id).
     */
    @Override
    public String getUsername() {
        return id;
    }
}