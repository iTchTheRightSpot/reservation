CREATE TABLE IF NOT EXISTS profile
(
    profile_id BIGSERIAL    NOT NULL UNIQUE PRIMARY KEY,
    firstname  VARCHAR(255) NOT NULL,
    lastname   VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL UNIQUE,
    image_key  VARCHAR(255)
);