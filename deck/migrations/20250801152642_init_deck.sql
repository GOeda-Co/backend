-- +goose Up
-- +goose StatementBegin

-- Deck Service Tables
CREATE TABLE IF NOT EXISTS decks (
    deck_id UUID PRIMARY KEY,
    created_by UUID NOT NULL,
    created_at TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_public BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_decks_created_by ON decks(created_by);
CREATE INDEX IF NOT EXISTS idx_decks_is_public ON decks(is_public);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS decks;

-- +goose StatementEnd
