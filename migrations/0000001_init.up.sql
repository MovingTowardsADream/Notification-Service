CREATE OR REPLACE FUNCTION make_uid() RETURNS text AS $$
DECLARE
new_uid text;
    done bool;
BEGIN
    done := false;
    WHILE NOT done LOOP
        new_uid := md5(''||now()::text||random()::text);
        done := NOT exists(SELECT 1 FROM users WHERE id=new_uid);
END LOOP;
RETURN new_uid;
END;
$$ LANGUAGE PLPGSQL VOLATILE;

CREATE TABLE IF NOT EXISTS users
(
    id TEXT DEFAULT make_uid()::text NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    phone VARCHAR(15) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS users_idx ON users(id);

CREATE TABLE IF NOT EXISTS notifications
(
    id BIGSERIAL PRIMARY KEY,
    email_notify BOOLEAN NOT NULL DEFAULT false,
    phone_notify BOOLEAN NOT NULL DEFAULT false,
    user_id TEXT NOT NULL REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS history_email_notify
(
    id BIGSERIAL PRIMARY KEY,
    request_id TEXT NOT NULL UNIQUE,
    subject TEXT,
    body TEXT,
    status TEXT CHECK (status IN ('processed', 'failed')),
    user_id TEXT NOT NULL REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS history_phone_notify
(
    id BIGSERIAL PRIMARY KEY,
    request_id TEXT NOT NULL UNIQUE,
    subject TEXT NULL,
    body TEXT NULL,
    status TEXT NOT NULL CHECK (status IN ('processed', 'failed')),
    user_id TEXT NOT NULL REFERENCES users(id)
);