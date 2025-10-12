package dev.whiteo.yadoma.repository;

import dev.whiteo.yadoma.exception.NotFoundObjectException;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.data.repository.NoRepositoryBean;
/**
 * Base interface for all MongoDB repositories in the application.
 * Provides common methods for entity retrieval and naming.
 * @param <E> entity type
 * @param <I> entity ID type
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@NoRepositoryBean
public interface AbstractRepository<E, I> extends MongoRepository<E, I> {

    /**
     * Finds an entity by its ID or throws NotFoundObjectException if not found.
     * @param id entity ID
     * @return entity instance
     */
    default E getOrThrow(I id) {
        return findById(id)
                .orElseThrow(() ->
                        new NotFoundObjectException(getEntityName() + " not found with id: " + id));
    }

    /**
     * Returns the simple name of the entity (without 'Repository' suffix).
     * @return entity name
     */
    default String getEntityName() {
        return getClass().getSimpleName().replace("Repository", "");
    }
}