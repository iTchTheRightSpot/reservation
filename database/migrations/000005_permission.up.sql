CREATE TYPE permissionenum AS ENUM ('READ', 'WRITE', 'DELETE');

CREATE TABLE IF NOT EXISTS permission
(
    permission_id    BIGSERIAL NOT NULL UNIQUE PRIMARY KEY,
    permission       permissionenum  NOT NULL,
    role_id    BIGINT NOT NULL,
    CONSTRAINT FK_permission_to_role_role_id
        FOREIGN KEY (role_id)
            REFERENCES role(role_id)
            ON DELETE CASCADE
            ON UPDATE RESTRICT
);