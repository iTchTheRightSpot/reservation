CREATE TABLE IF NOT EXISTS reservation_service
(
    junction_id    BIGSERIAL NOT NULL UNIQUE PRIMARY KEY,
    reservation_id BIGINT    NOT NULL,
    service_id     BIGINT    NOT NULL,
    CONSTRAINT FK_reservation_service_to_reservation_reservation_id
        FOREIGN KEY (reservation_id)
            REFERENCES reservation(reservation_id)
            ON DELETE RESTRICT
            ON UPDATE RESTRICT,
    CONSTRAINT FK_reservation_service_to_service_service_id
        FOREIGN KEY (service_id)
            REFERENCES service(service_id)
            ON DELETE RESTRICT
            ON UPDATE RESTRICT
);