DO $$

DECLARE
i integer := 0;
user_id text;

BEGIN
    WHILE i < 10 LOOP
        INSERT INTO users (username, email, phone, password_hash)
        VALUES (
               SUBSTRING(MD5(RANDOM()::text), 1, 10 + (RANDOM() * 10)::integer),
               SUBSTRING(MD5(RANDOM()::text), 1, 10 + (RANDOM() * 10)::integer) || '@gmail.com',
               LPAD(FLOOR(RANDOM() * (10^12 - 1) + 10^11)::text, 12, '0')::varchar(12),
               SUBSTRING(MD5(RANDOM()::text), 1, 10 + (RANDOM() * 10)::integer)
        ) ON CONFLICT DO NOTHING RETURNING id INTO user_id;

        INSERT INTO notifications (email_notify, phone_notify, user_id)
        VALUES (
               RANDOM() > 0.5,
               RANDOM() > 0.5,
               user_id
        ) ON CONFLICT DO NOTHING;
        i := i + 1;
END LOOP;
END;
$$;