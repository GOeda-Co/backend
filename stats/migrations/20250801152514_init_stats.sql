-- +goose Up
-- +goose StatementBegin

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
CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_deck_id ON reviews(deck_id);
CREATE INDEX IF NOT EXISTS idx_reviews_card_id ON reviews(card_id);
CREATE INDEX IF NOT EXISTS idx_reviews_created_at ON reviews(created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS reviews;

-- +goose StatementEnd
