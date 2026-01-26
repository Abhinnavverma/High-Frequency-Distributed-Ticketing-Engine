CREATE TABLE seats (
    id SERIAL PRIMARY KEY,
    row_number CHAR(1) NOT NULL,    -- Example: 'A', 'B', 'C'
    seat_number INT NOT NULL,       -- Example: 1, 2, 3
    status TEXT NOT NULL DEFAULT 'available', -- The crucial column
    price INT NOT NULL DEFAULT 50   -- Simplify: everything is $50 for now
);