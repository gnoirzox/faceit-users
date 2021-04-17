CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS country (
    alpha_code CHAR(3) NOT NULL,
    name VARCHAR NOT NULL,
    PRIMARY KEY (alpha_code)
);

CREATE TABLE IF NOT EXISTS users (
    _id uuid DEFAULT uuid_generate_v4 (),
    firstname VARCHAR NOT NULL,
    lastname VARCHAR NOT NULL,
    nickname VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    country_code VARCHAR NOT NULL,
    PRIMARY KEY (_id),
    CONSTRAINT fk_country FOREIGN KEY (country_code) REFERENCES country(alpha_code)
);

CREATE INDEX idx_users_nickname ON users(nickname);
