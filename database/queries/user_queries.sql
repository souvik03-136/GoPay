-- Create a new user
INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
VALUES (@id, @username, @email, @password_hash, NOW(), NOW())
RETURNING id;

-- Get user by ID
SELECT * FROM users
WHERE id = @id;

-- Get user by username
SELECT * FROM users
WHERE username = @username;

-- Get user by email
SELECT * FROM users
WHERE email = @email;

-- Update user
UPDATE users
SET username = @username, email = @email, password_hash = @password_hash, updated_at = NOW()
WHERE id = @id;

-- Delete user
DELETE FROM users
WHERE id = @id;
