CREATE TABLE reservations(
    id bigserial PRIMARY KEY,
    user_id bigint REFERENCES users(id),
    table_id bigint REFERENCES tables(id),
    seats_count integer NOT NULL,
    date date NOT NULL
);