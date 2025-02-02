CREATE TYPE paymentenum AS ENUM ('POS', 'ONLINE');

CREATE TABLE IF NOT EXISTS payment_detail
(
    payment_id     BIGSERIAL   NOT NULL UNIQUE PRIMARY KEY,
    reservation_id BIGSERIAL   NOT NULL UNIQUE,
    payment_type   paymentenum NOT NULL,
    transaction_id VARCHAR(255) UNIQUE DEFAULT NULL,
    is_paid        BOOLEAN             DEFAULT FALSE,
    CONSTRAINT FK_payment_detail_to_reservation_reservation_id
        FOREIGN KEY (reservation_id)
            REFERENCES reservation (reservation_id)
            ON DELETE RESTRICT
            ON UPDATE RESTRICT
);