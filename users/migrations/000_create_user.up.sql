CREATE TABLE app_user (
    id SERIAL PRIMARY KEY,
    username VARCHAR(16) NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    CONSTRAINT unique_username UNIQUE (username)
)