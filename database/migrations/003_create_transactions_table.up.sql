-- Create transactions table
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    payment_id INTEGER REFERENCES payments(id),
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL
);
