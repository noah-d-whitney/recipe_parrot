CREATE TABLE IF NOT EXISTS lists (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT now(),
    current BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS recipes (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    list_id BIGINT NOT NULL REFERENCES lists(id) ON DELETE CASCADE,
    title text NOT NULL,
    url text NOT NULL
);

CREATE TABLE IF NOT EXISTS ingredients (
    id bigserial PRIMARY KEY NOT NULL,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipe_id bigint NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    quantity text,
    unit text,
    name text
);

CREATE TABLE IF NOT EXISTS sites (
    name text PRIMARY KEY NOT NULL
);
