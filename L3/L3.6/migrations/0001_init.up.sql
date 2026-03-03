--items
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    type TEXT NOT NULL, -- income / expense
    amount NUMERIC(12,2) NOT NULL CHECK (amount >= 0),
    category TEXT,
    created_at TIMESTAMP NOT NULL
);