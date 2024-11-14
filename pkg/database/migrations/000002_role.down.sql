ALTER TABLE role
    DROP CONSTRAINT IF EXISTS FK_role_to_profile_profile_id;
DROP TABLE IF EXISTS role;
DROP TYPE IF EXISTS roleenum;
