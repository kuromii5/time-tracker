CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    passport_serie CHAR(4) NOT NULL,
    passport_number CHAR(6) NOT NULL,
    name VARCHAR(128) NOT NULL,
    surname VARCHAR(128) NOT NULL,
    patronymic VARCHAR(128),
    address TEXT NOT NULL,
    UNIQUE (passport_serie, passport_number)
);
CREATE INDEX idx_users_passport_serie ON users (passport_serie);
CREATE INDEX idx_users_passport_number ON users (passport_number);
CREATE INDEX idx_users_name ON users (name);
CREATE INDEX idx_users_surname ON users (surname);
CREATE INDEX idx_users_patronymic ON users (patronymic);
CREATE INDEX idx_users_address ON users (address);
CREATE INDEX idx_users_passport ON users (passport_serie, passport_number);
CREATE INDEX idx_users_created_at ON users (created_at);
CREATE INDEX idx_users_updated_at ON users (updated_at);

CREATE TABLE IF NOT EXISTS worklogs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    started_at TIMESTAMP DEFAULT NOW(),
    finished_at TIMESTAMP,
    task VARCHAR(255) NOT NULL,
    duration INTERVAL GENERATED ALWAYS AS (finished_at - started_at) STORED,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
CREATE INDEX idx_worklogs_user_id ON worklogs (user_id);
CREATE INDEX idx_worklogs_started_at ON worklogs (started_at);
CREATE INDEX idx_worklogs_finished_at ON worklogs (finished_at);