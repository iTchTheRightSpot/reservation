ALTER TABLE payment_detail
    DROP CONSTRAINT IF EXISTS FK_payment_detail_to_reservation_reservation_id;
DROP TABLE IF EXISTS payment_detail;
DROP TYPE IF EXISTS paymentenum;
