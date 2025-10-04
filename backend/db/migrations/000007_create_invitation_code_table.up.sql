ALTER TABLE servers
DROP COLUMN invite_code;

CREATE TABLE invitation_codes (
    id SERIAL PRIMARY KEY,
    server_id varchar(255) NOT NULL REFERENCES servers(server_id) ON DELETE CASCADE,
    token varchar(512) NOT NULL UNIQUE,
    status varchar(50) NOT NULL DEFAULT 'pending',
    expires_at bigint NOT NULL,
    created_at bigint NOT NULL,
    used_at bigint
);