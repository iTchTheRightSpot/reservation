CREATE TABLE IF NOT EXISTS schedule
(
    schedule_id     BIGSERIAL NOT NULL UNIQUE PRIMARY KEY,
    shift_start     TIMESTAMP NOT NULL,
    shift_end       TIMESTAMP NOT NULL,
    is_visible      BOOLEAN   NOT NULL DEFAULT FALSE,
    is_reoccurring  BOOLEAN   NOT NULL DEFAULT FALSE,
    staff_id        BIGINT    NOT NULL,
    CONSTRAINT FK_schedule_to_staff_staff_id
        FOREIGN KEY (staff_id)
            REFERENCES staff (staff_id)
            ON DELETE RESTRICT
            ON UPDATE RESTRICT
);

CREATE INDEX IX_schedule_composite1 ON schedule (staff_id, shift_start, shift_end);
CREATE INDEX IX_schedule_composite2 ON schedule (staff_id, shift_start, shift_end, is_visible);
CREATE INDEX IX_schedule_composite3 ON schedule (is_reoccurring, shift_start, shift_end);