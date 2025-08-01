-- +goose Up
-- +goose StatementBegin

-- SSO Service Tables
CREATE TABLE IF NOT EXISTS apps (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    secret VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    pass_hash BYTEA NOT NULL,
    name VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE
);

-- Deck Service Tables
CREATE TABLE IF NOT EXISTS decks (
    deck_id UUID PRIMARY KEY,
    created_by UUID NOT NULL,
    created_at TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_public BOOLEAN DEFAULT FALSE,
    card_quantity BIGINT,
    updated_at TIMESTAMP
);

-- Card Service Tables
CREATE TABLE IF NOT EXISTS cards (
    card_id UUID PRIMARY KEY,
    created_by UUID NOT NULL,
    created_at TIMESTAMP,
    word VARCHAR(100) NOT NULL,
    translation VARCHAR(100) NOT NULL,
    easiness DOUBLE PRECISION DEFAULT 2.5 NOT NULL,
    updated_at TIMESTAMP,
    interval SMALLINT DEFAULT 0,
    expires_at TIMESTAMP,
    repetition_number SMALLINT DEFAULT 0,
    deck_id UUID,
    tags TEXT[],
    is_public BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (deck_id) REFERENCES decks(deck_id) ON DELETE SET NULL
);

-- Stats Service Tables
CREATE TABLE IF NOT EXISTS reviews (
    result_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    deck_id UUID,
    card_id UUID,
    created_at TIMESTAMP,
    grade INTEGER NOT NULL CHECK (grade >= 0 AND grade <= 5)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_decks_created_by ON decks(created_by);
CREATE INDEX IF NOT EXISTS idx_decks_is_public ON decks(is_public);
CREATE INDEX IF NOT EXISTS idx_cards_created_by ON cards(created_by);
CREATE INDEX IF NOT EXISTS idx_cards_deck_id ON cards(deck_id);
CREATE INDEX IF NOT EXISTS idx_cards_expires_at ON cards(expires_at);
CREATE INDEX IF NOT EXISTS idx_cards_is_public ON cards(is_public);
CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_deck_id ON reviews(deck_id);
CREATE INDEX IF NOT EXISTS idx_reviews_card_id ON reviews(card_id);
CREATE INDEX IF NOT EXISTS idx_reviews_created_at ON reviews(created_at);

-- Insert default app for authentication
INSERT INTO apps (name, secret) 
VALUES ('repeatro', 'repeatro-secret-key') 
ON CONFLICT (id) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop tables in reverse order due to foreign key constraints
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS decks;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS apps;

-- +goose StatementEnd
