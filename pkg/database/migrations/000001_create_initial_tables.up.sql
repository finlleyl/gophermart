CREATE EXTENSION IF NOT EXISTS "pgcrypto";


CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       login TEXT UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE orders (
                        id SERIAL PRIMARY KEY,
                        user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                        number TEXT UNIQUE NOT NULL,
                        status TEXT CHECK (status IN ('NEW', 'PROCESSING', 'PROCESSED', 'INVALID')) DEFAULT 'NEW',
                        accrual NUMERIC(10, 2) DEFAULT 0,
                        uploaded_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE loyalty_accounts (
                                  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
                                  current_balance NUMERIC(10, 2) DEFAULT 0,
                                  withdrawn_balance NUMERIC(10, 2) DEFAULT 0
);

CREATE TABLE withdrawals (
                             id SERIAL PRIMARY KEY,
                             user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                             order_number TEXT NOT NULL,
                             sum NUMERIC(10, 2) NOT NULL,
                             withdrawn_at TIMESTAMP DEFAULT NOW()
);
