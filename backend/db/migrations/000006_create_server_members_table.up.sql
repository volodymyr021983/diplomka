CREATE TABLE server_members (
    id SERIAL PRIMARY KEY,
    user_id varchar(255) NOT NULL REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    server_id varchar(255) NOT NULL REFERENCES servers(server_id) ON DELETE CASCADE,
    user_role varchar(255) NOT NULL DEFAULT 'owner'
)