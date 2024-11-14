ALTER TABLE staff
    DROP CONSTRAINT IF EXISTS FK_staff_to_profile_profile_id;
DROP TABLE IF EXISTS staff;