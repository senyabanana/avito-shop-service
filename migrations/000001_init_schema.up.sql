CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    coins BIGINT NOT NULL DEFAULT 1000
);

CREATE TABLE IF NOT EXISTS transactions
(
    id SERIAL PRIMARY KEY,
    from_user INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS merch_items
(
    id SERIAL PRIMARY KEY,
    item_type VARCHAR(255) NOT NULL UNIQUE,
    price INT NOT NULL CHECK (price > 0)
);

CREATE TABLE IF NOT EXISTS inventory
(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    merch_id INT NOT NULL REFERENCES merch_items(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO merch_items (item_type, price) VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500);