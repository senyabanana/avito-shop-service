CREATE TABLE users
(
    id SERIAL NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    coins INT NOT NULL DEFAULT 1000
);

CREATE TABLE coins_transactions
(
    id SERIAL NOT NULL UNIQUE,
    from_user INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount INT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE inventory
(
    id SERIAL NOT NULL UNIQUE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_type VARCHAR(11) NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMP DEFAULT NOW()
);