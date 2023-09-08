CREATE TABLE currencies (
    code VARCHAR(5) UNIQUE NOT NULL
);

CREATE INDEX code_idx ON currencies (code);

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    currency_code INTEGER NOT NULL,
    number VARCHAR(90) NOT NULL,
    balance NUMERIC NOT NULL,
    FOREIGN KEY (currency_code) REFERENCES currencies (code),
    CHECK (balance >= 0)
);

CREATE INDEX number_idx ON accounts (number);

CREATE TABLE transactions (
    created_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(7) NOT NULL,
    amount NUMERIC NOT NULL,
    currency_code VARCHAR(5) NOT NULL,
    from_account INTEGER,
    to_account INTEGER NOT NULL,
    FOREIGN KEY (from_account) REFERENCES accounts (id),
    FOREIGN KEY (to_account) REFERENCES accounts (id)
);

