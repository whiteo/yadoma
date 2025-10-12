package dev.whiteo.yadoma.mapper;

import org.springframework.stereotype.Component;

/**
 * Base interface for mapping entities to DTO responses.
 * @param <E> entity type
 * @param <D> DTO type
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Component
public interface AbstractMapper<E, D> {

    /**
     * Converts an entity to its DTO response.
     * @param entity the entity to convert
     * @return DTO response
     */
    D toResponse(E entity);
}