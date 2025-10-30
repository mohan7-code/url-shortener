-- +goose Up
-- +goose StatementBegin
CREATE TABLE url_shortner (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT UNIQUE NOT NULL,
    click_count BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_accessed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_url_shortner_created_at ON url_shortner(created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS url_shortner;
-- +goose StatementEnd
