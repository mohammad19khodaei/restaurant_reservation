CREATE TABLE users(
    id bigserial PRIMARY KEY,
    username varchar UNIQUE NOT NULL,
    password varchar NOT NULL,
    created_at timestamp default now()
)