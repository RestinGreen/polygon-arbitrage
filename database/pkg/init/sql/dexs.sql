CREATE TABLE dexs (
    id SERIAL PRIMARY KEY,
    factory_address VARCHAR(255) NOT NULL UNIQUE,
    router_address VARCHAR(255) NOT NULL UNIQUE,
    num_pairs INTEGER
);


