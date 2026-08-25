CREATE TABLE
    IF NOT EXISTS rss_source (
        id BIGSERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        url TEXT NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS account (
        id BIGSERIAL PRIMARY KEY,
        login VARCHAR(255) NOT NULL UNIQUE,
        password_hash VARCHAR(60) NOT NULL
    );