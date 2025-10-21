-- 1) add column nullable (or with temporary default)
ALTER TABLE users ADD COLUMN email varchar(100);

-- 2) populate existing rows (choose a safe placeholder or compute from ID)
-- Example: generate unique placeholder emails using UUID-based id
UPDATE users
SET email = concat('user+', id::text, '@example.local')
WHERE email IS NULL;

-- 3) enforce NOT NULL (and remove temporary default if you used one)
ALTER TABLE users
ALTER COLUMN email SET NOT NULL;