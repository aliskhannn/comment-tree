package comment

import (
	"context"

	"github.com/google/uuid"

	"github.com/aliskhannn/comment-tree/internal/model"
)

// Repository provides methods for interacting with the comments table.
type Repository interface {
	CreateComment(ctx context.Context, comment *model.Comment) (uuid.UUID, error)
	GetCommentsByParentID(ctx context.Context, parentID uuid.UUID) ([]model.Comment, error)
	DeleteComment(ctx context.Context, id uuid.UUID) error
}

// Service provides methods for interacting with the comments table.
type Service struct {
	repo Repository
}

// NewService creates a new Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateComment creates a new comment.
func (s *Service) CreateComment(ctx context.Context, comment *model.Comment) (uuid.UUID, error) {
	return s.repo.CreateComment(ctx, comment)
}

// GetCommentsByParentID returns the comment with the given ID and all nested descendants.
func (s *Service) GetCommentsByParentID(ctx context.Context, parentID uuid.UUID) ([]model.Comment, error) {
	return s.repo.GetCommentsByParentID(ctx, parentID)
}

// DeleteComment deletes a comment by ID and all nested descendants.
func (s *Service) DeleteComment(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteComment(ctx, id)
}