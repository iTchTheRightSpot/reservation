CREATE TABLE IF NOT EXISTS shift
(
    shift_id       BIGSERIAL NOT NULL UNIQUE,
    start          TIMESTAMP NOT NULL,
    shift_end      TIMESTAMP NOT NULL,
    is_enabled     BOOLEAN   NOT NULL DEFAULT FALSE,
    is_reoccurring BOOLEAN   NOT NULL DEFAULT FALSE,
    staff_id       BIGINT    NOT NULL,
    CONSTRAINT FK_shift_to_staff_staff_id
        FOREIGN KEY (staff_id)
            REFERENCES staff (staff_id)
            ON DELETE RESTRICT
            ON UPDATE RESTRICT
);

CREATE INDEX IX_shift_composite1 ON shift (staff_id, start, shift_end);
CREATE INDEX IX_shift_composite2 ON shift (staff_id, start, shift_end, is_enabled);
CREATE INDEX IX_shift_composite3 ON shift (is_reoccurring, start, shift_end);