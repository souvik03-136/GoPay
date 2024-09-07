-- Create a new payment
INSERT INTO payments (id, user_id, amount, currency, status, created_at)
VALUES (@id, @user_id, @amount, @currency, @status, NOW());

-- Get payment by ID
SELECT * FROM payments
WHERE id = @id;

-- List all payments for a user
SELECT * FROM payments
WHERE user_id = @user_id;

-- Update payment status
UPDATE payments
SET status = @status
WHERE id = @id;
