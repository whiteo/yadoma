package dev.whiteo.yadoma.util;

import lombok.experimental.UtilityClass;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;

/**
 * Utility class for password hashing and verification using BCrypt.
 * Provides static methods for encoding and matching passwords.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@UtilityClass
public class PasswordUtil {

    private static final PasswordEncoder encoder = new BCryptPasswordEncoder();

    /**
     * Hashes a raw password using BCrypt.
     * @param rawPassword plain text password
     * @return hashed password string
     */
    public static String hash(String rawPassword) {
        return encoder.encode(rawPassword);
    }

    /**
     * Verifies that a raw password matches the given hash.
     * @param rawPassword plain text password
     * @param hash hashed password
     * @return true if the password matches, false otherwise
     */
    public static boolean matches(String rawPassword, String hash) {
        return encoder.matches(rawPassword, hash);
    }
}