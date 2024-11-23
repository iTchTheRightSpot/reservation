ALTER TABLE staff_service
    DROP CONSTRAINT IF EXISTS FK_staff_service_to_service_service_id;

ALTER TABLE staff_service
    DROP CONSTRAINT IF EXISTS FK_staff_service_to_staff_staff_id;

DROP TABLE IF EXISTS staff_service;