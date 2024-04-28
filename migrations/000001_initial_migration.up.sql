CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    first_name text,
    last_name text,
    phone_number text NOT NULL
);
