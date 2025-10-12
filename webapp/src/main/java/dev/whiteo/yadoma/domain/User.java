package dev.whiteo.yadoma.domain;

import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.data.mongodb.core.index.Indexed;
import org.springframework.data.mongodb.core.mapping.Document;
import org.springframework.data.mongodb.core.mapping.Field;

import java.util.List;

/**
 * Represents a user entity stored in MongoDB.
 * Contains user email, password hash, and role.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Getter
@Setter
@NoArgsConstructor
@Document(collection = "users")
public class User extends AbstractDocument {

    /**
     * User's email address. Must be unique.
     */
    @Field("email")
    @Indexed(unique = true)
    private String email;

    /**
     * Hashed password of the user.
     */
    @Field("password_hash")
    private String passwordHash;

    /**
     * List of container IDs associated with the user.
     */
    @Field("containers")
    private List<String> containerIds;
    /**
     * Role of the user in the application.
     */
    @Field("role")
    private Role role = Role.USER;
}