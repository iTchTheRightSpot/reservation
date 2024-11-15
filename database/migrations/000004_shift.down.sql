ALTER TABLE shift
    DROP CONSTRAINT IF EXISTS FK_shift_to_staff_staff_id;

ALTER TABLE shift
    DROP CONSTRAINT IF EXISTS IX_shift_composite1;

ALTER TABLE shift
    DROP CONSTRAINT IF EXISTS IX_shift_composite2;

DROP TABLE IF EXISTS shift;