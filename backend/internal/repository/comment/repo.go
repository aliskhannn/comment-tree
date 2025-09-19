package comment

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/dbpg"

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
func (r *Repository) CreateComment(ctx context.Context, comment *model.Comment) (uuid.UUID, error) {
	query := `
		INSERT INTO comments (id, parent_id, content)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id uuid.UUID
	err := r.db.Master.QueryRowContext(
		ctx, query,
		comment.ID, comment.ParentID, comment.Content,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return id, nil
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
			COALESCE(parent_id, '00000000-0000-0000-0000-000000000000'::uuid) AS parent_id,
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

	comments := []model.Comment{}
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

	if len(comments) == 0 {
		return nil, ErrCommentNotFound
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

	rows, err := r.db.Master.ExecContext(ctx, query, id)
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
