package dev.whiteo.yadoma.util;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class PasswordUtilTest {

    @Test
    void hash_ShouldReturnHashedPassword() {
        String rawPassword = "testPassword123";
        String hashedPassword = PasswordUtil.hash(rawPassword);

        assertNotNull(hashedPassword);
        assertNotEquals(rawPassword, hashedPassword);
        assertTrue(hashedPassword.startsWith("$2a$"));
        assertTrue(hashedPassword.length() > 50);
    }

    @Test
    void hash_WithDifferentPasswords_ShouldProduceDifferentHashes() {
        String password1 = "password1";
        String password2 = "password2";

        String hash1 = PasswordUtil.hash(password1);
        String hash2 = PasswordUtil.hash(password2);

        assertNotEquals(hash1, hash2);
    }

    @Test
    void hash_WithSamePassword_ShouldProduceDifferentHashesDueToSalt() {
        String password = "testPassword";

        String hash1 = PasswordUtil.hash(password);
        String hash2 = PasswordUtil.hash(password);

        assertNotEquals(hash1, hash2);
    }

    @Test
    void matches_WithCorrectPassword_ShouldReturnTrue() {
        String rawPassword = "testPassword123";
        String hashedPassword = PasswordUtil.hash(rawPassword);

        boolean matches = PasswordUtil.matches(rawPassword, hashedPassword);

        assertTrue(matches);
    }

    @Test
    void matches_WithIncorrectPassword_ShouldReturnFalse() {
        String rawPassword = "testPassword123";
        String wrongPassword = "wrongPassword";
        String hashedPassword = PasswordUtil.hash(rawPassword);

        boolean matches = PasswordUtil.matches(wrongPassword, hashedPassword);

        assertFalse(matches);
    }

    @Test
    void matches_WithEmptyPassword_ShouldHandleCorrectly() {
        String emptyPassword = "";
        String hashedPassword = PasswordUtil.hash(emptyPassword);

        boolean matches = PasswordUtil.matches(emptyPassword, hashedPassword);

        assertTrue(matches);
    }

    @Test
    void matches_WithNullPassword_ShouldHandleCorrectly() {
        String hashedPassword = PasswordUtil.hash("somePassword");

        assertThrows(IllegalArgumentException.class, () ->
            PasswordUtil.matches(null, hashedPassword));
    }
}