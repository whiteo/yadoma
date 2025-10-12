package dev.whiteo.yadoma.mapper;

import org.junit.jupiter.api.Test;
import org.springframework.stereotype.Component;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Test class for AbstractMapper interface.
 */
class AbstractMapperTest {

    private static class TestEntity {
        private String name;

        public TestEntity(String name) {
            this.name = name;
        }

        public String getName() {
            return name;
        }
    }

    private static class TestDto {
        private String name;

        public TestDto(String name) {
            this.name = name;
        }

        public String getName() {
            return name;
        }
    }

    @Component
    private static class TestMapper implements AbstractMapper<TestEntity, TestDto> {
        @Override
        public TestDto toResponse(TestEntity entity) {
            if (entity == null) {
                return null;
            }
            return new TestDto(entity.getName());
        }
    }

    @Test
    void interface_ShouldHaveComponentAnnotation() {
        assertTrue(AbstractMapper.class.isAnnotationPresent(Component.class));
    }

    @Test
    void implementation_ShouldMapEntityToDto() {
        TestMapper mapper = new TestMapper();
        TestEntity entity = new TestEntity("Test Name");

        TestDto dto = mapper.toResponse(entity);

        assertNotNull(dto);
        assertEquals("Test Name", dto.getName());
    }

    @Test
    void implementation_ShouldHandleNullEntity() {
        TestMapper mapper = new TestMapper();

        TestDto dto = mapper.toResponse(null);

        assertNull(dto);
    }
}
