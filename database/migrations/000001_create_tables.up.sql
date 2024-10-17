-- Install golang migrate: go get -u -d github.com/golang-migrate/migrate/cmd/migrate
-- Run migrations: migrate -path=./database/migrations -database "postgres://postgres:postgres@localhost:5432/matcha?sslmode=disable" up

-- Create riders table
CREATE TABLE IF NOT EXISTS riders (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- Create drivers table
CREATE TABLE IF NOT EXISTS drivers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    geohash VARCHAR(12),
    status VARCHAR(20) DEFAULT 'available' -- 'available', 'on_trip'
);

-- Create trips table
CREATE TABLE IF NOT EXISTS trips (
    id SERIAL PRIMARY KEY,
    rider_id INT REFERENCES riders(id),
    driver_id INT REFERENCES drivers(id),
    start_latitude DOUBLE PRECISION,
    start_longitude DOUBLE PRECISION,
    end_latitude DOUBLE PRECISION,
    end_longitude DOUBLE PRECISION,
    status VARCHAR(20) DEFAULT 'requested' -- 'requested', 'accepted', 'completed'
);
