-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_comments_parent_id;
-- +goose StatementEnd
