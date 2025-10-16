package comment

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"

	"github.com/aliskhannn/comment-tree/internal/model"
)

var ErrCommentNotFound = errors.New("comment not found")

// Repository provides methods for interacting with the comments table.
type Repository struct {
	db *dbpg.DB
}

// NewRepository creates a new Repository.
func NewRepository(db *dbpg.DB) *Repository {
	return &Repository{db: db}
}

// CreateComment creates a new comment.
func (r *Repository) CreateComment(ctx context.Context, comment *model.Comment) (model.Comment, error) {
	query := `
		INSERT INTO comments (parent_id, content)
		VALUES ($1, $2)
		RETURNING id, parent_id, content, created_at, updated_at
	`

	zlog.Logger.Printf("repo: parent id: %v", comment.ParentID)

	var c model.Comment

	err := r.db.QueryRowContext(
		ctx, query,
		comment.ParentID, comment.Content,
	).Scan(&c.ID, &c.ParentID, &c.Content, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return model.Comment{}, fmt.Errorf("failed to create comment: %w", err)
	}

	return c, nil
}

// GetCommentsByParentID returns the comment with the given ID and all nested descendants.
func (r *Repository) GetCommentsByParentID(ctx context.Context, parentID uuid.UUID) ([]model.Comment, error) {
	query := `
		WITH RECURSIVE comment_tree AS (
			SELECT id, parent_id, content, created_at, updated_at
			FROM comments
			WHERE id = $1
			UNION ALL
			SELECT c.id, c.parent_id, c.content, c.created_at, c.updated_at
			FROM comments c
			JOIN comment_tree ct ON c.parent_id = ct.id
		)
		SELECT 
			id,
			parent_id,
			content,
			created_at,
			updated_at
		FROM comment_tree
		ORDER BY created_at
	`

	rows, err := r.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by parent ID: %w", err)
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(&comment.ID, &comment.ParentID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get comments by parent ID: %w", err)
	}

	return comments, nil
}

// GetComments retrieves comments by parent ID with optional search, sorting, and pagination.
func (r *Repository) GetComments(ctx context.Context, parentID *uuid.UUID, search string, sort string, limit, offset int) ([]model.Comment, error) {
	query := `SELECT id, parent_id, content, created_at, updated_at FROM comments WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if parentID != nil {
		query += fmt.Sprintf(" AND parent_id = $%d", argIdx)
		args = append(args, *parentID)
		argIdx++
	}

	if search != "" {
		query += fmt.Sprintf(" AND to_tsvector('english', content) @@ plainto_tsquery('english', $%d)", argIdx)
		args = append(args, search)
		argIdx++
	}

	allowedSorts := map[string]string{
		"created_asc":  "created_at ASC",
		"created_desc": "created_at DESC",
		"updated_asc":  "updated_at ASC",
		"updated_desc": "updated_at DESC",
	}

	if sortSQL, ok := allowedSorts[sort]; ok {
		query += " ORDER BY " + sortSQL
	} else {
		query += " ORDER BY created_at DESC" // По умолчанию
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.ParentID, &c.Content, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	return comments, nil
}

// DeleteComment deletes a comment by ID and all nested descendants.
func (r *Repository) DeleteComment(ctx context.Context, id uuid.UUID) error {
	query := `
		WITH RECURSIVE to_delete AS (
    		SELECT id FROM comments WHERE id = $1
    		UNION ALL
    		SELECT c.id
    		FROM comments c
    		JOIN to_delete td ON c.parent_id = td.id
		)
		DELETE FROM comments
		WHERE id IN (SELECT id FROM to_delete);
	`

	rows, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	n, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if n == 0 {
		return ErrCommentNotFound
	}

	return nil
}
