CREATE TABLE IF NOT EXISTS staff
(
    staff_id BIGSERIAL NOT NULL UNIQUE ,
    staff_uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    bio VARCHAR(255),
    profile_id BIGINT,
    CONSTRAINT FK_staff_to_profile_profile_id
        FOREIGN KEY (profile_id)
        REFERENCES profile(profile_id)
        ON DELETE SET NULL
        ON UPDATE RESTRICT
);