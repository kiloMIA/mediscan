CREATE TABLE IF NOT EXISTS appointments (
    id SERIAL PRIMARY KEY,
    discord_id VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    details TEXT
);

CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY,
    discord_id VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL
);

