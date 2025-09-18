-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_comments_created_at ON comments(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_comments_created_at;
-- +goose StatementEnd
