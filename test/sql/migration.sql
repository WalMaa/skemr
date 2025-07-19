

-- 1. Rename column 'username' to 'user_name'
ALTER TABLE users
    RENAME COLUMN username TO user_name;

-- 2. Drop the column 'password_hash'
ALTER TABLE users
    DROP COLUMN password_hash;

-- 3. Add a new column 'is_active' with default value
ALTER TABLE users
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;


ALTER TABLE users
    ALTER COLUMN email TYPE TEXT;

-- 5. Rename table 'users' to 'app_users'
ALTER TABLE users
    RENAME TO app_users;

-- 6. Add an index on 'user_name'
CREATE INDEX IF NOT EXISTS idx_app_users_user_name ON app_users(user_name);
