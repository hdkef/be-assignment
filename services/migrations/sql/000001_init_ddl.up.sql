-- create schemas
CREATE SCHEMA accounts;
CREATE SCHEMA transactions;

-- Create 'users' table
CREATE TABLE accounts.users (
    id UUID PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    dob DATE,
    job VARCHAR(50) NOT NULL,
    created_at TIMESTAMP
);

-- Create 'user_address' table
CREATE TABLE accounts.user_address (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES accounts.users(id),
    zip INTEGER,
    address TEXT NOT NULL,
    district VARCHAR(50) NOT NULL,
    city VARCHAR(50) NOT NULL,
    province VARCHAR(50) NOT NULL,
    country VARCHAR(50) NOT NULL
);

-- Create 'accounts' table
CREATE TABLE accounts.accounts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    user_id UUID REFERENCES accounts.users(id),
    acc_type VARCHAR(10) CHECK (acc_type IN ('DEBIT', 'CREDIT', 'LOAN')),
    acc_desc VARCHAR(200),
    currency VARCHAR(5),
    balance DECIMAL(10, 2)
);

-- Create 'histories' table
CREATE TABLE accounts.histories (
    id UUID PRIMARY KEY,
    reff_num UUID,
    created_at TIMESTAMP,
    acc_id UUID REFERENCES accounts.accounts(id) NOT NULL,
    trx_type VARCHAR(10) CHECK (trx_type IN ('POSITIVE', 'NEGATIVE')),
    amount DECIMAL(10, 2),
    description VARCHAR(200),
    status VARCHAR(10) CHECK (status IN ('FAILED', 'SUCCESS')),
    acc_id_2 UUID REFERENCES accounts.accounts(id)
);

-- Create 'accounts_balance' table
CREATE TABLE transactions.accounts_balance (
    user_id UUID NOT NULL,
    acc_id UUID NOT NULL UNIQUE,
    currency VARCHAR(5),
    balance DECIMAL(10, 2),
    PRIMARY KEY (acc_id)
);

-- Create 'transactions_logs' table
CREATE TABLE transactions.transactions_logs (
    reff_num UUID,
    created_at TIMESTAMP,
    acc_id UUID REFERENCES transactions.accounts_balance(acc_id),
    amount DECIMAL(10, 2),
    description VARCHAR(200),
    status VARCHAR(10) CHECK (status IN ('FAILED', 'SUCCESS')),
    event_type VARCHAR(10) CHECK (event_type IN ('SEND', 'WITHDRAW')),
    to_acc_id UUID REFERENCES transactions.accounts_balance(acc_id)
);

-- Create 'scheduled_trx' table
CREATE TABLE transactions.scheduled_trx (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    status VARCHAR(10) CHECK (status IN ('ACTIVE', 'INACTIVE')),
    acc_id UUID REFERENCES transactions.accounts_balance(acc_id),
    type VARCHAR(10) CHECK (type IN ('SEND', 'WITHDRAW')),
    amount DECIMAL(10, 2),
    schedule VARCHAR(10),
    to_acc_id UUID,
    has_checked BOOLEAN,
    last_checked DATE,
);

-- Create 'queued_trx' table
CREATE TABLE transactions.queued_trx (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    status VARCHAR(10) CHECK (status IN ('PENDING', 'EXECUTED','')),
    result VARCHAR(10) CHECK (result IN ('FAILED', 'SUCCESS')),
    schedule_trx_id UUID REFERENCES transactions.scheduled_trx(id)
);