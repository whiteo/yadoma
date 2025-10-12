package dev.whiteo.yadoma.service;

import dev.whiteo.yadoma.domain.Role;
import dev.whiteo.yadoma.domain.User;
import dev.whiteo.yadoma.dto.token.TokenResponse;
import dev.whiteo.yadoma.dto.token.TokenValidationResponse;
import dev.whiteo.yadoma.dto.user.UserCreateRequest;
import dev.whiteo.yadoma.dto.user.UserLoginRequest;
import dev.whiteo.yadoma.dto.user.UserResponse;
import dev.whiteo.yadoma.dto.user.UserUpdatePasswordRequest;
import dev.whiteo.yadoma.exception.ExecutionConflictException;
import dev.whiteo.yadoma.mapper.UserMapper;
import dev.whiteo.yadoma.repository.UserRepository;
import dev.whiteo.yadoma.security.TokenInteract;
import dev.whiteo.yadoma.security.UserDetailsImpl;
import dev.whiteo.yadoma.util.PasswordUtil;
import jakarta.servlet.http.HttpServletRequest;
import lombok.RequiredArgsConstructor;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.stream.Collectors;

/**
 * Service class for user-related operations such as authentication, token management, and user lookup.
 * Implements Spring Security's UserDetailsService for authentication integration.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Service
@RequiredArgsConstructor
public class UserService implements UserDetailsService {

    private final TokenInteract tokenInteract;
    private final UserRepository repository;
    private final UserMapper mapper;

    /**
     * Retrieves the details of all users.
     *
     * @param userId ID of the user making the request
     * @return ResponseEntity containing the list of user details or error response
     */
    public List<UserResponse> findAll(String userId) {
        User user = repository.getOrThrow(userId);
        if (user.getRole() != Role.ADMIN) {
            throw new BadCredentialsException("Access denied");
        }

        return repository.findAll().stream()
                .map(mapper::toResponse)
                .collect(Collectors.toList());
    }

    /**
     * Finds a user by ID and returns a UserResponse DTO.
     *
     * @param userId ID of the user to find
     * @return UserResponse containing user details
     */
    public UserResponse me(String userId) {
        return mapper.toResponse(repository.getOrThrow(userId));
    }

    /**
     * Creates a new user after checking for unique email and hashing the password.
     *
     * @param request request containing user creation details
     */
    public void create(UserCreateRequest request) {
        checkUniqueEmail(request.email());
        repository.save(mapper.toEntity(request.email(), PasswordUtil.hash(request.password())));
    }

    /**
     * Removes a user by ID, ensuring the user has permission to delete.
     *
     * @param deleteId ID of the user to delete
     * @param userId   ID of the user performing the deletion
     */
    public void remove(String deleteId, String userId) {
        repository.delete(validateUserAccess(deleteId, userId));
    }

    /**
     * Updates the password of an existing user after verifying the old password.
     *
     * @param userId  ID of the user to update
     * @param request request containing old and new passwords
     */
    public void updatePassword(String userId, UserUpdatePasswordRequest request) {
        User user = repository.getOrThrow(userId);
        verifyPassword(request.oldPassword(), user.getPasswordHash());
        user.setPasswordHash(PasswordUtil.hash(request.password()));
    }

    /**
     * Authenticates a user by ID or email and returns a login response with JWT token.
     *
     * @param request user login request
     * @return login response containing JWT token and user info
     */
    public TokenResponse getToken(UserLoginRequest request) {
        User user = findByEmail(request.email());
        verifyPassword(request.password(), user.getPasswordHash());
        return new TokenResponse(tokenInteract.generateToken(loadUserByUsername(user.getId())));
    }

    /**
     * Validates the JWT token from the HTTP request.
     *
     * @param request HTTP servlet request
     * @return TokenValidationResponse if token is valid
     */
    public TokenValidationResponse validateToken(HttpServletRequest request) {
        String token = tokenInteract.getToken(request);
        return new TokenValidationResponse(tokenInteract.validateToken(token));
    }

    /**
     * Loads user details by user ID for authentication purposes.
     *
     * @param userId user's ID
     * @return UserDetails implementation containing user info
     * @throws UsernameNotFoundException if user is not found
     */
    @Override
    public UserDetails loadUserByUsername(String userId) throws UsernameNotFoundException {
        User user = repository.getOrThrow(userId);
        return UserDetailsImpl.builder()
                .id(user.getId())
                .build();
    }

    /**
     * Checks if a user owns a container by ID.
     *
     * @param containerId container's ID
     * @param userId user's ID
     * @return true if user owns the container, false otherwise
     */
    public Boolean isContainerIdContains(String containerId, String userId) {
        User user = repository.getOrThrow(userId);
        return user.getContainerIds().contains(containerId) || user.getRole() == Role.ADMIN;
    }

    /**
     * Validates user access to a container.
     *
     * @param userId user's ID
     * @param adminId admin's ID
     * @return user object if access is granted, throws exception otherwise
     */
    public User validateUserAccess(String userId, String adminId) {
        User user = repository.getOrThrow(userId);
        if (!userId.equals(adminId) && user.getRole() != Role.ADMIN) {
            throw new BadCredentialsException("Access denied");
        }
        return user;
    }

    /**
     * Retrieves a user by ID.
     *
     * @param userId user's ID
     * @return user object
     */
    public User getUserById(String userId) {
        return repository.getOrThrow(userId);
    }

    /**
     * Removes a container ID from a user's list of container IDs and saves the user.
     *
     * @param containerId ID of the container to remove
     * @param user        user object from which to remove the container ID
     */
    public void removeContainerFromUser(String containerId, User user) {
        user.getContainerIds().remove(containerId);
        repository.save(user);
    }

    /**
     * Adds a container ID to a user's list of container IDs and saves the user.
     *
     * @param containerId ID of the container to add
     * @param user        user object to which to add the container ID
     */
    public void addContainerToUser(String containerId, User user) {
        user.getContainerIds().add(containerId);
        repository.save(user);
    }

    /**
     * Finds a user by email or throws BadCredentialsException if not found.
     *
     * @param email user's email address
     * @return User entity
     */
    private User findByEmail(String email) {
        return repository.findByEmailIgnoreCase(email)
                .orElseThrow(() -> new BadCredentialsException(
                        "User not found with email: " + email
                ));
    }

    /**
     * Verifies the raw password against the stored password hash.
     * Throws BadCredentialsException if the password is invalid.
     *
     * @param rawPassword  plain text password
     * @param passwordHash hashed password
     */
    private void verifyPassword(String rawPassword, String passwordHash) {
        if (!PasswordUtil.matches(rawPassword, passwordHash)) {
            throw new BadCredentialsException("Invalid password");
        }
    }

    /**
     * Checks if the email is unique in the database.
     * Throws ExecutionConflictException if a user with the email already exists.
     *
     * @param email user's email address
     */
    private void checkUniqueEmail(String email) {
        repository.findByEmailIgnoreCase(email)
                .ifPresent(u -> {
                    throw new ExecutionConflictException(
                            "User with email '" + email + "' already exists");
                });
    }
}