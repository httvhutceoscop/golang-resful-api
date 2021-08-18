CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS base_table (
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE user_account (
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
    name VARCHAR(50) UNIQUE NOT NULL,
    initial_balance INTEGER NOT NULL,
    description VARCHAR(255),
    currency_id uuid,
    user_account_id uuid,
    FOREIGN KEY (currency_id) REFERENCES currency (id) ON DELETE SET NULL,
    FOREIGN KEY (user_account_id) REFERENCES user_account (id) ON DELETE CASCADE
) INHERITS (base_table);

CREATE TABLE transaction (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
    amount INTEGER NOT NULL,
    description VARCHAR(255),
    account_id uuid,
    transaction_type_id uuid,
    FOREIGN KEY (account_id) REFERENCES account (id) ON DELETE SET NULL,
    FOREIGN KEY (transaction_type_id) REFERENCES transaction_type (id) ON DELETE SET NULL
) INHERITS (base_table);