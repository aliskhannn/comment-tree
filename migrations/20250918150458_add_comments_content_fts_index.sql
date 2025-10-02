-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_comments_content_fts ON comments USING GIN (to_tsvector('english', content));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_comments_content_fts;
-- +goose StatementEnd
