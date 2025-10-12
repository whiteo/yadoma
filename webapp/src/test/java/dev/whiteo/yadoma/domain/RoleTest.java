package dev.whiteo.yadoma.domain;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class RoleTest {

    @Test
    void testRoleValues() {
        Role[] roles = Role.values();
        assertEquals(2, roles.length);
        assertTrue(java.util.Arrays.asList(roles).contains(Role.USER));
        assertTrue(java.util.Arrays.asList(roles).contains(Role.ADMIN));
    }

    @Test
    void testRoleValueOf() {
        assertEquals(Role.USER, Role.valueOf("USER"));
        assertEquals(Role.ADMIN, Role.valueOf("ADMIN"));
    }

    @Test
    void testRoleValueOfInvalid() {
        assertThrows(IllegalArgumentException.class, () -> Role.valueOf("INVALID"));
    }

    @Test
    void testRoleOrdinal() {
        assertEquals(0, Role.ADMIN.ordinal());
        assertEquals(1, Role.USER.ordinal());
    }

    @Test
    void testRoleName() {
        assertEquals("USER", Role.USER.name());
        assertEquals("ADMIN", Role.ADMIN.name());
    }
}
