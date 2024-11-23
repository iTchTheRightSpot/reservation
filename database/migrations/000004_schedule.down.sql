ALTER TABLE schedule
    DROP CONSTRAINT IF EXISTS FK_schedule_to_staff_staff_id;

ALTER TABLE schedule
    DROP CONSTRAINT IF EXISTS IX_schedule_composite1;

ALTER TABLE schedule
    DROP CONSTRAINT IF EXISTS IX_schedule_composite2;

ALTER TABLE schedule
    DROP CONSTRAINT IF EXISTS IX_schedule_composite3;

DROP TABLE IF EXISTS schedule;