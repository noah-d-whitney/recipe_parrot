CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY NOT NULL,
    first_name text,
    last_name text,
    phone_number text NOT NULL
);

CREATE TABLE IF NOT EXISTS shopping_trips (
    id bigserial PRIMARY KEY NOT NULL,
    user_id bigint REFERENCES users(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS meals (
    id bigserial PRIMARY KEY NOT NULL,
    user_id bigint REFERENCES users(id) NOT NULL, 
    plan bigint REFERENCES shopping_trips(id) NOT NULL,
    url text NOT NULL,
    name text NOT NULL,
    minutes int,
    ingredients json
);
