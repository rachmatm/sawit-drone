-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

CREATE TABLE estates (
    id UUID PRIMARY KEY,
    width INTEGER NOT NULL,
    length INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CHECK(width BETWEEN 1 AND 50000),
    CHECK(length BETWEEN 1 AND 50000)
);

CREATE TABLE trees (
    id UUID PRIMARY KEY,
    estate_id UUID NOT NULL,
    x INTEGER NOT NULL,
    y INTEGER NOT NULL,
    height INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_estate
        FOREIGN KEY (estate_id)
        REFERENCES estates(id)
        ON DELETE CASCADE,

    CHECK(height BETWEEN 1 AND 30)
);

CREATE UNIQUE INDEX idx_unique_tree_plot
ON trees(estate_id, x, y);

CREATE INDEX idx_tree_estate
ON trees(estate_id);