ALTER TABLE reservation_service
    DROP CONSTRAINT IF EXISTS FK_reservation_service_to_reservation_reservation_id;

ALTER TABLE reservation_service
    DROP CONSTRAINT IF EXISTS FK_reservation_service_to_service_service_id;

DROP TABLE IF EXISTS reservation_service