-- Send a new email
INSERT INTO emails (id, user_id, subject, body, sent_at)
VALUES (@id, @user_id, @subject, @body, NOW());

-- Get email by ID
SELECT * FROM emails
WHERE id = @id;

-- List all emails for a user
SELECT * FROM emails
WHERE user_id = @user_id;
