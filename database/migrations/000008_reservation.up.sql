CREATE TYPE reservationenum AS ENUM ('CONFIRMED', 'CANCELLED');

CREATE TABLE IF NOT EXISTS reservation
(
    reservation_id BIGSERIAL       NOT NULL UNIQUE PRIMARY KEY,
    name           VARCHAR(100)    NOT NULL,
    email          VARCHAR(320)    NOT NULL,
    description    VARCHAR(255) DEFAULT NULL,
    address        VARCHAR(255) DEFAULT NULL,
    phone          VARCHAR(20)  DEFAULT NULL,
    image_key      VARCHAR(37)  DEFAULT NULL,
    price          DECIMAL(6, 2)   NOT NULL,
    status         reservationenum NOT NULL,
    created_at     TIMESTAMP       NOT NULL,
    scheduled_for  TIMESTAMP       NOT NULL,
    expire_at      TIMESTAMP       NOT NULL,
    staff_id       BIGINT          NOT NULL,
    CONSTRAINT FK_reservation_to_staff_staff_id
        FOREIGN KEY (staff_id)
            REFERENCES staff(staff_id)
            ON DELETE RESTRICT
            ON UPDATE RESTRICT
);

CREATE INDEX IX_reservation_email ON reservation (email);
CREATE INDEX IX_reservation_composite1 ON reservation (staff_id, scheduled_for, expire_at, status);
CREATE INDEX IX_reservation_composite2 ON reservation (email, scheduled_for, expire_at, status);