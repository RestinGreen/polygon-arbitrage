CREATE TABLE tokens ( 
    id SERIAL PRIMARY KEY, 
    address VARCHAR(255) NOT NULL UNIQUE, symbol VARCHAR(255), 
    name VARCHAR(255)
);
