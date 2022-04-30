CREATE TYPE currency AS ENUM ('Euro', 'Dollar');

CREATE TABLE IF NOT EXISTS beer(
    id VARCHAR(255) UNIQUE NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    name VARCHAR(255) NOT NULL,
    price_currency currency NOT NULL,
    price_amount INT NOT NULL CONSTRAINT positive_price CHECK (price_amount > 0)
);