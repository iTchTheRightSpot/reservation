CREATE TABLE IF NOT EXISTS staff_service
(
    junction_id BIGSERIAL NOT NULL UNIQUE PRIMARY KEY,
    staff_id    BIGINT    NOT NULL,
    service_id  BIGINT    NOT NULL,
    CONSTRAINT FK_staff_service_to_service_service_id
        FOREIGN KEY (service_id)
            REFERENCES service(service_id)
            ON DELETE RESTRICT
            ON UPDATE RESTRICT,
    CONSTRAINT FK_staff_service_to_staff_staff_id
        FOREIGN KEY (staff_id)
            REFERENCES staff(staff_id)
            ON DELETE RESTRICT
            ON UPDATE RESTRICT
);