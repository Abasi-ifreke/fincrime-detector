-- init.sql

CREATE TABLE IF NOT EXISTS detection_rules (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    priority INT,
    condition TEXT NOT NULL,
    weight FLOAT
);

CREATE TABLE IF NOT EXISTS transactions (
    id TEXT PRIMARY KEY,
    account_id TEXT NOT NULL,
    transaction_type TEXT,
    amount FLOAT,
    timestamp TIMESTAMP
);

CREATE TABLE IF NOT EXISTS alerts (
    id TEXT PRIMARY KEY,
    transaction_id TEXT,
    account_id TEXT,
    reason TEXT,
    score FLOAT,
    timestamp TIMESTAMP
);

-- Example rules
INSERT INTO detection_rules (name, description, condition, weight) VALUES
('Large Transaction', 'Transaction amount exceeds a threshold', 'large_amount', 0.4),
('Blacklisted Account', 'Transaction from a known blacklisted account', 'blacklisted_account', 0.8);