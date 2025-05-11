CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(255) PRIMARY KEY,
    account_id VARCHAR(255) NOT NULL,
    transaction_type VARCHAR(50) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS detection_rules (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    priority INTEGER NOT NULL DEFAULT 1,
    condition VARCHAR(255) NOT NULL,
    weight DECIMAL(5, 2) NOT NULL DEFAULT 0.1
);

-- Example rules
INSERT INTO detection_rules (name, description, condition, weight) VALUES
('Large Transaction', 'Transaction amount exceeds a threshold', 'large_amount', 0.4),
('Blacklisted Account', 'Transaction from a known blacklisted account', 'blacklisted_account', 0.8);