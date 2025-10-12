package dev.whiteo.yadoma.domain;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.time.LocalDateTime;

import static org.junit.jupiter.api.Assertions.*;

class AbstractDocumentTest {

    private AbstractDocument document;

    @BeforeEach
    void setUp() {
        document = new AbstractDocument() {};
    }

    @Test
    void testIdSetterAndGetter() {
        String id = "test-id-123";
        document.setId(id);
        assertEquals(id, document.getId());
    }

    @Test
    void testCreationDateSetterAndGetter() {
        LocalDateTime creationDate = LocalDateTime.now();
        document.setCreationDate(creationDate);
        assertEquals(creationDate, document.getCreationDate());
    }

    @Test
    void testModifyDateSetterAndGetter() {
        LocalDateTime modifyDate = LocalDateTime.now();
        document.setModifyDate(modifyDate);
        assertEquals(modifyDate, document.getModifyDate());
    }

    @Test
    void testInitialValues() {
        AbstractDocument newDocument = new AbstractDocument() {};
        assertNull(newDocument.getId());
        assertNull(newDocument.getCreationDate());
        assertNull(newDocument.getModifyDate());
    }

    @Test
    void testTimestampOrder() {
        LocalDateTime creationDate = LocalDateTime.now();
        LocalDateTime modifyDate = creationDate.plusMinutes(5);

        document.setCreationDate(creationDate);
        document.setModifyDate(modifyDate);

        assertTrue(document.getModifyDate().isAfter(document.getCreationDate()));
    }
}
