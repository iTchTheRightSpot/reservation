ALTER TABLE permission
    DROP CONSTRAINT IF EXISTS FK_permission_to_role_role_id;
DROP TABLE IF EXISTS permission;
DROP TYPE IF EXISTS permissionenum;
