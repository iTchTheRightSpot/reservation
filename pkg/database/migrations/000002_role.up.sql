CREATE TYPE roleenum AS ENUM ('STAFF', 'DEVELOPER', 'USER');

CREATE TABLE IF NOT EXISTS role
(
    role_id    BIGSERIAL NOT NULL UNIQUE,
    role       roleenum  NOT NULL DEFAULT ('USER'),
    profile_id BIGINT    NOT NULL,
    CONSTRAINT FK_role_to_profile_profile_id
        FOREIGN KEY (profile_id)
            REFERENCES profile(profile_id)
            ON DELETE CASCADE
            ON UPDATE RESTRICT
);