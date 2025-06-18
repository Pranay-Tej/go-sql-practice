-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS projects(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE
) 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE projects;
-- +goose StatementEnd