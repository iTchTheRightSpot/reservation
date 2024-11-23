ALTER TABLE reservation
    DROP CONSTRAINT IF EXISTS FK_reservation_to_staff_staff_id;

DROP INDEX IF EXISTS IX_reservation_email;
DROP INDEX IF EXISTS IX_reservation_composite1;
DROP INDEX IF EXISTS IX_reservation_composite2;
DROP TABLE IF EXISTS reservation;
DROP TYPE IF EXISTS reservationenum;