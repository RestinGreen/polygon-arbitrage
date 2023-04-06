# How to set up

1. Have a polygon node ipc
2. 
    - Create pgsql user "bot" with a password.
    - CREATE USER bot WITH ENCRYPTED PASSWORD '';
    - Create database: CREATE DATABASE arbitrage
    - log in to bot user with psql -U bot -d arbitrage -W.
     If failed, than find the "pg_hba.conf" file (find / -name pg_hba.conf) and add the line "local | all | bot | (leave empty) | md5.


# Generate bindings like this

abigen  --abi univ2pair.abi --pkg univ2pair --type UniV2Pair --out univ2pair.go

# Database structure

CREATE TABLE dexs (
    id SERIAL PRIMARY KEY,
    factory_address VARCHAR(255) NOT NULL UNIQUE,
    router_address VARCHAR(255) NOT NULL UNIQUE,
    num_pairs INTEGER
);


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
    FOREIGN KEY (dex_id) REFERENCES dexs(id)
);

