-- CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE IF NOT EXISTS schedule
(
    schedule_id    BIGSERIAL NOT NULL UNIQUE PRIMARY KEY,
    schedule_start TIMESTAMP NOT NULL,
    schedule_end   TIMESTAMP NOT NULL,
    is_visible     BOOLEAN   NOT NULL DEFAULT FALSE,
    is_reoccurring BOOLEAN   NOT NULL DEFAULT FALSE,
    staff_id       BIGINT    NOT NULL,
    CONSTRAINT FK_schedule_to_staff_staff_id
        FOREIGN KEY (staff_id)
            REFERENCES staff (staff_id)
            ON DELETE RESTRICT
            ON UPDATE RESTRICT,
    CONSTRAINT EX_schedule_overlap_constraint
        EXCLUDE USING gist (
            staff_id WITH =, tstzrange(schedule_start, schedule_end) WITH &&
        )
);

CREATE INDEX IX_schedule_composite1 ON schedule (staff_id, schedule_start, schedule_end);
CREATE INDEX IX_schedule_composite2 ON schedule (staff_id, schedule_start, schedule_end, is_visible);
CREATE INDEX IX_schedule_composite3 ON schedule (schedule_start, schedule_end, is_visible, is_reoccurring);