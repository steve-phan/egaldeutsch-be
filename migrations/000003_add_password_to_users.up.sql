-- Add password column to users table if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'users' AND column_name = 'password') THEN
        -- 1) add the column (nullable)
        ALTER TABLE users ADD COLUMN password varchar(255);

        -- 2) backfill existing rows with a random placeholder
        UPDATE users
        SET password = gen_random_uuid()::text
        WHERE password IS NULL;

        -- 3) enforce NOT NULL constraint
        ALTER TABLE users
        ALTER COLUMN password SET NOT NULL;
    END IF;
END $$;
