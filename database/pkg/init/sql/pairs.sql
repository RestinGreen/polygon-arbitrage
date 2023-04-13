CREATE TABLE pairs (
    id SERIAL PRIMARY KEY,
    pair_address VARCHAR(255) NOT NULL UNIQUE,
    token0_address VARCHAR(255) NOT NULL,
    token1_address VARCHAR(255) NOT NULL,
    reserve0 NUMERIC NOT NULL,
    reserve1 NUMERIC NOT NULL,
    last_updated INTEGER,
    dex_id INTEGER,
    used INTEGER DEFAULT 0,
    token0_id INTEGER,
    token1_id INTEGER,
    FOREIGN KEY (token0_id) REFERENCES tokens(id),
    FOREIGN KEY (token1_id) REFERENCES tokens(id),
    FOREIGN KEY (dex_id) REFERENCES dexs(id)
);
