-- +goose Up
-- +goose StatementBegin

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
    is_public BOOLEAN DEFAULT FALSE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_cards_created_by ON cards(created_by);
CREATE INDEX IF NOT EXISTS idx_cards_deck_id ON cards(deck_id);
CREATE INDEX IF NOT EXISTS idx_cards_expires_at ON cards(expires_at);
CREATE INDEX IF NOT EXISTS idx_cards_is_public ON cards(is_public);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS cards;

-- +goose StatementEnd
