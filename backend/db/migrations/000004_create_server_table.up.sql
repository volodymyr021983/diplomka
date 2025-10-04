CREATE TABLE servers (
    id SERIAL PRIMARY KEY,
    server_id varchar(255) UNIQUE NOT NULL,
    servername varchar(255) NOT NULL,
    owner_id varchar(255) NOT NULL REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    invite_code varchar(255) UNIQUE NULL
)