CREATE TABLE currencies (
    code VARCHAR(5) UNIQUE NOT NULL
);

CREATE INDEX code_idx ON currencies (code);

CREATE TABLE accounts (
    number VARCHAR(42) UNIQUE NOT NULL,
    currency_code VARCHAR(5) NOT NULL,
    balance NUMERIC NOT NULL,
    FOREIGN KEY (currency_code) REFERENCES currencies (code),
    CHECK (balance >= 0)
);

CREATE INDEX number_idx ON accounts (number);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(7) NOT NULL,
    amount NUMERIC NOT NULL,
    currency_code VARCHAR(5) NOT NULL,
    from_account VARCHAR(90),
    to_account VARCHAR(90) NOT NULL
);

INSERT INTO currencies VALUES ('RUB');
INSERT INTO accounts VALUES ('1','RUB', 0);
