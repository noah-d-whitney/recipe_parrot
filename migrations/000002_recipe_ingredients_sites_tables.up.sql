CREATE TABLE IF NOT EXISTS recipes (
    id bigserial PRIMARY KEY NOT NULL,
    user_id bigint REFERENCES users(id) NOT NULL,
    title text NOT NULL
);

CREATE TABLE IF NOT EXISTS ingredients (
    id bigserial PRIMARY KEY NOT NULL,
    user_id bigint REFERENCES users(id) NOT NULL,
    recipe_id bigint REFERENCES recipes(id) NOT NULL,
    quantity text,
    unit text,
    name text
);

CREATE TABLE IF NOT EXISTS sites (
    name text PRIMARY KEY NOT NULL
);
