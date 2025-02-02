CREATE TABLE IF NOT EXISTS profile
(
    profile_id BIGSERIAL    NOT NULL UNIQUE PRIMARY KEY,
    firstname  VARCHAR(50)  NOT NULL,
    lastname   VARCHAR(50)  NOT NULL,
    email      VARCHAR(320) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    locked     BOOLEAN DEFAULT FALSE,
    image_key  VARCHAR(255)
);