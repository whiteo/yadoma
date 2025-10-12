package dev.whiteo.yadoma.domain;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.Id;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.mongodb.core.mapping.Field;

import java.time.LocalDateTime;

/**
 * Abstract base class for all MongoDB entities.
 * Contains common fields such as id, creation date, and modification date.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Getter
@Setter
public abstract class AbstractDocument {

    /**
     * Unique identifier of the document.
     */
    @Id
    private String id;

    /**
     * Date and time when the document was last modified.
     */
    @LastModifiedDate
    @Field("modify_date")
    private LocalDateTime modifyDate;

    /**
     * Date and time when the document was created.
     */
    @CreatedDate
    @Field("creation_date")
    private LocalDateTime creationDate;
}