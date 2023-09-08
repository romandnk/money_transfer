CREATE TABLE currencies (
    id SERIAL PRIMARY KEY,
    code VARCHAR(5) UNIQUE NOT NULL
);

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    currency_id INTEGER NOT NULL,
    balance NUMERIC NOT NULL,
    FOREIGN KEY (currency_id) REFERENCES currencies (id),
    CHECK (balance >= 0)
);

CREATE TABLE transactions (
    status VARCHAR(7) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,

);