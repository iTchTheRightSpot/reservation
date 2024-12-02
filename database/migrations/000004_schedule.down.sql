ALTER TABLE schedule DROP CONSTRAINT IF EXISTS FK_schedule_to_staff_staff_id;

ALTER TABLE schedule DROP CONSTRAINT IF EXISTS EX_schedule_overlap_constraint;

DROP INDEX IF EXISTS IX_schedule_composite1;

DROP INDEX IF EXISTS IX_schedule_composite2;

DROP INDEX IF EXISTS IX_schedule_composite3;

DROP TABLE IF EXISTS schedule;