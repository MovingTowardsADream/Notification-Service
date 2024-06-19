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
    username TEXT not null unique,
    email TEXT not null unique,
    phone varchar(11) not null unique,
    password_hash TEXT not null,
    time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS notifications
(
    id          serial       not null unique,
    email_notify BOOLEAN NOT NULL,
    phone_notify BOOLEAN NOT NULL,
    user_id TEXT NOT NULL REFERENCES users(id)
);