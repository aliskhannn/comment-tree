-- +goose Up
-- +goose StatementBegin
ALTER TABLE comments ADD CONSTRAINT fk_parent_id FOREIGN KEY (parent_id) REFERENCES comments(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE comments DROP CONSTRAINT fk_parent_id;
-- +goose StatementEnd
