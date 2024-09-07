-- Create a new transaction
INSERT INTO transactions (id, user_id, amount, transaction_date, status)
VALUES (@id, @user_id, @amount, @transaction_date, @status);

-- Get transaction by ID
SELECT * FROM transactions
WHERE id = @id;

-- Get transactions by user ID
SELECT * FROM transactions
WHERE user_id = @user_id;

-- Update transaction
UPDATE transactions
SET amount = @amount, transaction_date = @transaction_date, status = @status
WHERE id = @id;

-- Delete transaction
DELETE FROM transactions
WHERE id = @id;
