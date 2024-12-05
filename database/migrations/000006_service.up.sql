CREATE TABLE IF NOT EXISTS service
(
    service_id      BIGSERIAL       NOT NULL UNIQUE PRIMARY KEY,
    name            VARCHAR(50)     NOT NULL UNIQUE,
    price           DECIMAL(6, 2)   NOT NULL,
    is_visible      BOOL            NOT NULL DEFAULT FALSE,
    is_reoccurring  BOOL            NOT NULL DEFAULT FALSE,
    duration        INTEGER         NOT NULL,
    clean_up_time   INTEGER         NOT NULL
);