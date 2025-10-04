CREATE TABLE channels (
    id SERIAL PRIMARY KEY,
    channel_id varchar(255) UNIQUE NOT NULL,
    own_server_id varchar(255) NOT NULL REFERENCES servers(server_id) ON DELETE CASCADE,
    channel_type varchar(255) NOT NULL,
    channelname varchar(255) NOT NULL
)