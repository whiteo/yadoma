package dev.whiteo.yadoma.repository;

import dev.whiteo.yadoma.exception.NotFoundObjectException;
import org.junit.jupiter.api.Test;
import org.springframework.lang.NonNull;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;

class AbstractRepositoryTest {

    interface SampleEntity {}

    static class SampleRepository implements AbstractRepository<SampleEntity, String> {
        @Override public <S extends SampleEntity> @NonNull S save(@NonNull S entity) { return entity; }
        @Override public <S extends SampleEntity> @NonNull java.util.List<S> saveAll(@NonNull Iterable<S> entities) { java.util.ArrayList<S> list = new java.util.ArrayList<>(); entities.forEach(list::add); return list; }
        @Override public @NonNull Optional<SampleEntity> findById(@NonNull String s) { return Optional.empty(); }
        @Override public boolean existsById(@NonNull String s) { return false; }
        @Override public @NonNull java.util.List<SampleEntity> findAll() { return new java.util.ArrayList<>(); }
        @Override public @NonNull java.util.List<SampleEntity> findAllById(@NonNull Iterable<String> strings) { return new java.util.ArrayList<>(); }
        @Override public long count() { return 0; }
        @Override public void deleteById(@NonNull String s) {}
        @Override public void delete(@NonNull SampleEntity entity) {}
        @Override public void deleteAllById(@NonNull Iterable<? extends String> strings) {}
        @Override public void deleteAll(@NonNull Iterable<? extends SampleEntity> entities) {}
        @Override public void deleteAll() {}
        @Override
        public <S extends SampleEntity> @NonNull S insert(@NonNull S entity) {
            return entity;
        }
        @Override
        public <S extends SampleEntity> @NonNull java.util.List<S> insert(@NonNull Iterable<S> entities) {
            java.util.ArrayList<S> list = new java.util.ArrayList<>();
            entities.forEach(list::add);
            return list;
        }
        @Override
        public @NonNull java.util.List<SampleEntity> findAll(@NonNull org.springframework.data.domain.Sort sort) {
            return new java.util.ArrayList<>();
        }
        @Override
        public @NonNull org.springframework.data.domain.Page<SampleEntity> findAll(@NonNull org.springframework.data.domain.Pageable pageable) {
            return org.springframework.data.domain.Page.empty();
        }
        @Override
        public <S extends SampleEntity> @NonNull java.util.List<S> findAll(@NonNull org.springframework.data.domain.Example<S> example) { return new java.util.ArrayList<>(); }
        @Override
        public <S extends SampleEntity> @NonNull java.util.List<S> findAll(@NonNull org.springframework.data.domain.Example<S> example, @NonNull org.springframework.data.domain.Sort sort) { return new java.util.ArrayList<>(); }
        @Override
        public <S extends SampleEntity> @NonNull org.springframework.data.domain.Page<S> findAll(@NonNull org.springframework.data.domain.Example<S> example, @NonNull org.springframework.data.domain.Pageable pageable) {
            return org.springframework.data.domain.Page.empty();
        }
        @Override
        public <S extends SampleEntity> @NonNull Optional<S> findOne(@NonNull org.springframework.data.domain.Example<S> example) { return Optional.empty(); }
        @Override
        public <S extends SampleEntity> long count(@NonNull org.springframework.data.domain.Example<S> example) { return 0; }
        @Override
        public <S extends SampleEntity> boolean exists(@NonNull org.springframework.data.domain.Example<S> example) { return false; }
        @Override
        public <S extends SampleEntity, R> @NonNull R findBy(
            @NonNull org.springframework.data.domain.Example<S> example,
            @NonNull java.util.function.Function<org.springframework.data.repository.query.FluentQuery.FetchableFluentQuery<S>, R> queryFunction
        ) {
            return java.util.Objects.requireNonNull(queryFunction.apply(null));
        }
    }

    @Test
    void getEntityName_shouldRemoveRepositorySuffix() {
        SampleRepository repo = new SampleRepository();
        String name = repo.getEntityName();
        assertEquals("Sample", name);
    }

    @Test
    void getOrThrow_whenNotFound_throwsNotFound() {
        SampleRepository repo = new SampleRepository();

        NotFoundObjectException ex = assertThrows(NotFoundObjectException.class, () -> repo.getOrThrow("1"));
        assertTrue(ex.getMessage().contains("Sample not found"));
    }

    @Test
    void getOrThrow_throwsNotFoundObjectException() {
        SampleRepository repo = new SampleRepository();
        assertThrows(NotFoundObjectException.class, () -> repo.getOrThrow("id"));
    }

    @Test
    void getEntityName_returnsCorrectName() {
        SampleRepository repo = new SampleRepository();
        assertEquals("Sample", repo.getEntityName());
    }
}
