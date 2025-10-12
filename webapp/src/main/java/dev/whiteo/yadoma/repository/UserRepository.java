package dev.whiteo.yadoma.repository;

import dev.whiteo.yadoma.domain.User;
import org.springframework.stereotype.Repository;

import java.util.Optional;

/**
 * Repository interface for User entities.
 * Provides methods for user-specific queries.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Repository
public interface UserRepository extends AbstractRepository<User, String>{

    /**
     * Finds a user by email, ignoring case.
     * @param email user's email address
     * @return Optional containing the user if found, or empty otherwise
     */
    Optional<User> findByEmailIgnoreCase(String email);
}