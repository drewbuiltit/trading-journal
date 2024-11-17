CREATE TABLE trades
(
    id         SERIAL PRIMARY KEY,
    user_id    INT         NOT NULL REFERENCES users (id),
    symbol     VARCHAR(50) NOT NULL,
    quantity   DECIMAL     NOT NULL,
    price      DECIMAL     NOT NULL,
    trade_date TIMESTAMP DEFAULT NOW(),
    strategy   VARCHAR(255),
    note       VARCHAR(1000),
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_user_id ON trades (user_id);