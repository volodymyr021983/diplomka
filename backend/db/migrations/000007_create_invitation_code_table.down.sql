DROP TABLE IF EXISTS invitation_codes;

ALTER TABLE servers
ADD COLUMN invite_code varchar(255);