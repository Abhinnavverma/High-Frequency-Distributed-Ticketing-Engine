CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    seat_id INT NOT NULL REFERENCES UNIQUE seats(id), -- Foreign Key! Links to seats table
    user_id INT NOT NULL,                      -- Just an ID for now
    created_at TIMESTAMP DEFAULT NOW()
);