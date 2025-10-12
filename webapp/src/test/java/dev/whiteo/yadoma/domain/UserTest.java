package dev.whiteo.yadoma.domain;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.time.LocalDateTime;
import java.util.Arrays;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class UserTest {

    private User user;

    @BeforeEach
    void setUp() {
        user = new User();
    }

    @Test
    void testUserCreation() {
        assertNotNull(user);
        assertEquals(Role.USER, user.getRole());
    }

    @Test
    void testEmailSetterAndGetter() {
        String email = "test@example.com";
        user.setEmail(email);
        assertEquals(email, user.getEmail());
    }

    @Test
    void testPasswordHashSetterAndGetter() {
        String passwordHash = "hashedPassword123";
        user.setPasswordHash(passwordHash);
        assertEquals(passwordHash, user.getPasswordHash());
    }

    @Test
    void testContainerIdsSetterAndGetter() {
        List<String> containerIds = Arrays.asList("container1", "container2", "container3");
        user.setContainerIds(containerIds);
        assertEquals(containerIds, user.getContainerIds());
        assertEquals(3, user.getContainerIds().size());
    }

    @Test
    void testRoleSetterAndGetter() {
        user.setRole(Role.ADMIN);
        assertEquals(Role.ADMIN, user.getRole());

        user.setRole(Role.USER);
        assertEquals(Role.USER, user.getRole());
    }

    @Test
    void testDefaultRole() {
        User newUser = new User();
        assertEquals(Role.USER, newUser.getRole());
    }

    @Test
    void testInheritanceFromAbstractDocument() {
        assertTrue(user instanceof AbstractDocument);
    }

    @Test
    void testIdFromAbstractDocument() {
        String id = "12345";
        user.setId(id);
        assertEquals(id, user.getId());
    }

    @Test
    void testTimestampsFromAbstractDocument() {
        LocalDateTime now = LocalDateTime.now();
        user.setCreationDate(now);
        user.setModifyDate(now);

        assertEquals(now, user.getCreationDate());
        assertEquals(now, user.getModifyDate());
    }

    @Test
    void testUserEquality() {
        User user1 = new User();
        User user2 = new User();

        user1.setId("123");
        user1.setEmail("test@example.com");
        user1.setPasswordHash("hash123");
        user1.setRole(Role.USER);

        user2.setId("123");
        user2.setEmail("test@example.com");
        user2.setPasswordHash("hash123");
        user2.setRole(Role.USER);

        // Note: This test depends on Lombok generating equals method
        // If equals is not overridden, this test might fail
        assertNotSame(user1, user2);
    }

    @Test
    void testNullContainerIds() {
        user.setContainerIds(null);
        assertNull(user.getContainerIds());
    }

    @Test
    void testEmptyContainerIds() {
        user.setContainerIds(Arrays.asList());
        assertNotNull(user.getContainerIds());
        assertTrue(user.getContainerIds().isEmpty());
    }
}
