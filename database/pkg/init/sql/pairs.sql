CREATE TABLE pairs (
    id SERIAL PRIMARY KEY,
    pair_address VARCHAR(255) NOT NULL UNIQUE,
    reserve0 NUMERIC NOT NULL,
    reserve1 NUMERIC NOT NULL,
    last_updated INTEGER,
    dex_id INTEGER,
    token0_id INTEGER,
    token1_id INTEGER,
    FOREIGN KEY (token0_id) REFERENCES tokens(id),
    FOREIGN KEY (token1_id) REFERENCES tokens(id),
    FOREIGN KEY (dex_id) REFERENCES dexs(id)
);
