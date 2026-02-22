-- USERS
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    telegram_id TEXT,
    role VARCHAR(20) NOT NULL DEFAULT 'user'
        CHECK (role IN ('user', 'admin')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- EVENTS
CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    total_seats INT NOT NULL CHECK (total_seats > 0),
    available_seats INT NOT NULL
        CHECK (available_seats >= 0 AND available_seats <= total_seats),
    booking_ttl INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- BOOKINGS
CREATE TABLE bookings (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id),
    status TEXT NOT NULL CHECK (status IN ('pending', 'confirmed', 'canceled')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

-- =========================
-- INDEXES (ТОЛЬКО ПОСЛЕ CREATE TABLE bookings)
-- =========================

-- Индекс для поиска просроченных броней
CREATE INDEX idx_bookings_expires_at
    ON bookings (expires_at)
    WHERE status = 'pending';

-- Индекс для выборки броней по событию
CREATE INDEX idx_bookings_event_id
    ON bookings(event_id);

-- Запрет дублирующих активных броней одного пользователя на одно событие
CREATE UNIQUE INDEX idx_unique_user_event
    ON bookings(user_id, event_id)
    WHERE status IN ('pending', 'confirmed');

-- Индекс для выборки броней пользователя
CREATE INDEX idx_bookings_user_id
    ON bookings(user_id);
