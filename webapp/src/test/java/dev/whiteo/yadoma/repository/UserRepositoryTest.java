package dev.whiteo.yadoma.repository;

import dev.whiteo.yadoma.domain.Role;
import dev.whiteo.yadoma.domain.User;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.data.mongo.DataMongoTest;
import org.springframework.test.context.ActiveProfiles;

import java.util.List;
import java.util.Optional;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

@ActiveProfiles("test")
@DataMongoTest
class UserRepositoryTest {

    @Autowired
    private UserRepository userRepository;

    private User testUser;

    @BeforeEach
    void setUp() {
        userRepository.deleteAll();

        testUser = new User();
        testUser.setEmail("test@example.com");
        testUser.setPasswordHash("hashedPassword");
        testUser.setContainerIds(List.of("container1", "container2"));
        testUser.setRole(Role.USER);

        userRepository.save(testUser);
    }

    @Test
    void findByEmailIgnoreCase_ShouldReturnUser() {
        Optional<User> result = userRepository.findByEmailIgnoreCase("test@example.com");

        assertTrue(result.isPresent());
        assertEquals("test@example.com", result.get().getEmail());
        assertEquals("hashedPassword", result.get().getPasswordHash());
    }

    @Test
    void findByEmailIgnoreCase_WithDifferentCase_ShouldReturnUser() {
        Optional<User> result = userRepository.findByEmailIgnoreCase("TEST@EXAMPLE.COM");

        assertTrue(result.isPresent());
        assertEquals("test@example.com", result.get().getEmail());
    }

    @Test
    void findByEmailIgnoreCase_WithNonExistentEmail_ShouldReturnEmpty() {
        Optional<User> result = userRepository.findByEmailIgnoreCase("nonexistent@example.com");

        assertFalse(result.isPresent());
    }

    @Test
    void saveAndFindById_ShouldWork() {
        User newUser = new User();
        newUser.setEmail("new@example.com");
        newUser.setPasswordHash("newPassword");
        newUser.setRole(Role.USER);

        User savedUser = userRepository.save(newUser);
        Optional<User> foundUser = userRepository.findById(savedUser.getId());

        assertNotNull(savedUser.getId());
        assertTrue(foundUser.isPresent());
        assertEquals("new@example.com", foundUser.get().getEmail());
    }

    @Test
    void deleteById_ShouldRemoveUser() {
        String userId = testUser.getId();

        userRepository.deleteById(userId);
        boolean exists = userRepository.existsById(userId);

        assertFalse(exists);
    }

    @Test
    void findAll_ShouldReturnAllUsers() {
        List<User> all = userRepository.findAll();

        assertEquals(1, all.size());
    }

    @Test
    void count_ShouldReturnCorrectCount() {
        long count = userRepository.count();

        assertEquals(1, count);
    }
}