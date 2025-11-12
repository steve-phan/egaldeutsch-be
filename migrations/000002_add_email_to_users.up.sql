-- Add email column to users table if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'users' AND column_name = 'email') THEN
        -- 1) add column nullable
        ALTER TABLE users ADD COLUMN email varchar(100);

        -- 2) populate existing rows
        UPDATE users
        SET email = concat('user+', id::text, '@example.local')
        WHERE email IS NULL;

        -- 3) enforce NOT NULL
        ALTER TABLE users
        ALTER COLUMN email SET NOT NULL;
    END IF;
END $$;