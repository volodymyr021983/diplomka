CREATE TABLE ip_addresses (
    id SERIAL PRIMARY KEY,
    ip_address VARCHAR(255) NOT NULL,
    email_verify_time bigint NULL,
    reset_pwd_time bigint NULL
);