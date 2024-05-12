CREATE TABLE IF NOT EXISTS permissions (
                                           id bigserial PRIMARY KEY,
                                           code text NOT NULL
);
CREATE TABLE IF NOT EXISTS user_info_permissions (
                                                 user_id bigint NOT NULL REFERENCES user_info ON DELETE CASCADE,
                                                 permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
                                                 PRIMARY KEY (user_id, permission_id)
);-- Add the two permissions to the table.
INSERT INTO permissions (code)
VALUES
    ('user'),
    ('admin');