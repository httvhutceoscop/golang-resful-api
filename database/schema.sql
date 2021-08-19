CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS base_table (
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE auth_user (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
) INHERITS (base_table);

CREATE TABLE transaction_type (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
    code VARCHAR(10) UNIQUE NOT NULL,
    description VARCHAR(255)
);

CREATE TABLE currency (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
    code VARCHAR(10) UNIQUE NOT NULL,
    description VARCHAR(255)
);

CREATE TABLE account (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
    name VARCHAR(50) NOT NULL,
    initial_balance INTEGER NOT NULL,
    description VARCHAR(255),
    currency_id uuid NOT NULL,
    auth_user_id uuid NOT NULL,
    FOREIGN KEY (currency_id) REFERENCES currency (id) ON DELETE SET NULL,
    FOREIGN KEY (auth_user_id) REFERENCES auth_user (id) ON DELETE CASCADE,
    UNIQUE (name, auth_user_id)
) INHERITS (base_table);
CREATE INDEX index_name ON account(name);

CREATE TABLE transaction (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
    amount INTEGER NOT NULL,
    description VARCHAR(255),
    account_id uuid NOT NULL,
    transaction_type_id uuid NOT NULL,
    FOREIGN KEY (account_id) REFERENCES account (id) ON DELETE SET NULL,
    FOREIGN KEY (transaction_type_id) REFERENCES transaction_type (id) ON DELETE SET NULL
) INHERITS (base_table);
CREATE INDEX index_created_at ON transaction(created_at DESC);

INSERT INTO transaction_type (code, description) VALUES('W', 'Withdrawals');
INSERT INTO transaction_type (code, description) VALUES('D', 'Deposits');

INSERT INTO currency (code, description) VALUES('VND', 'Vietnamese Dong');
INSERT INTO currency (code, description) VALUES('JPY', 'Japanese Yen');
INSERT INTO currency (code, description) VALUES('THB', 'Thai Baht');