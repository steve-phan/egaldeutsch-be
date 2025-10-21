-- Add password column to users in a safe way:
-- 1) add the column (nullable)
-- 2) backfill existing rows with a random placeholder (gen_random_uuid())
-- 3) set the column to NOT NULL

ALTER TABLE users ADD COLUMN password varchar(255);

-- Backfill existing rows with a non-guessable placeholder so we can enforce NOT NULL
UPDATE users
SET password = gen_random_uuid()::text
WHERE password IS NULL;

-- Now enforce NOT NULL constraint
ALTER TABLE users
ALTER COLUMN password SET NOT NULL;
