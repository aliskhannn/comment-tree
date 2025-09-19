package comment

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"

	"github.com/aliskhannn/comment-tree/internal/api/respond"
	"github.com/aliskhannn/comment-tree/internal/model"
	"github.com/aliskhannn/comment-tree/internal/repository/comment"
)

// Service is the interface for the comment service.
type Service interface {
	CreateComment(ctx context.Context, comment *model.Comment) (uuid.UUID, error)
	GetCommentsByParentID(ctx context.Context, parentID uuid.UUID) ([]model.Comment, error)
	DeleteComment(ctx context.Context, id uuid.UUID) error
}

// Handler is the handler for the comment API.
type Handler struct {
	service Service
}

// NewHandler creates a new Handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// CreateRequest is the request for the create comment API.
type CreateRequest struct {
	ParentID uuid.UUID `json:"parent_id"`
	Content  string    `json:"content" binding:"required,min=1,max=1000"`
}

// Create creates a new comment.
func (h *Handler) Create(c *ginext.Context) {
	// Bind the request.
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to bind JSON")
		respond.Fail(c.Writer, http.StatusBadRequest, err)
		return
	}

	// Create the comment.
	cm := &model.Comment{
		ParentID: req.ParentID,
		Content:  req.Content,
	}

	id, err := h.service.CreateComment(c, cm)
	if err != nil {
		respond.Fail(c.Writer, http.StatusInternalServerError, err)
		return
	}

	respond.Created(c.Writer, id)
}

// Get retrieves the comment with the given ID and all nested descendants.
func (h *Handler) Get(c *ginext.Context) {
	// Extract parentID from query string (?parent=...).
	parentID := c.Query("parent")
	if parentID == "" {
		zlog.Logger.Warn().Msg("parent id is required")
		respond.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("parent id is required"))
		return
	}

	// Get comments.
	comments, err := h.service.GetCommentsByParentID(c.Request.Context(), uuid.MustParse(parentID))
	if err != nil {
		// If no comments are found, return 404 Not Found.
		if errors.Is(err, comment.ErrCommentNotFound) {
			zlog.Logger.Error().Err(err).Msg("comment not found")
			respond.Fail(c.Writer, http.StatusNotFound, fmt.Errorf("comment not found"))
		}

		zlog.Logger.Error().Err(err).Msg("failed to get comments")
		respond.Fail(c.Writer, http.StatusInternalServerError, fmt.Errorf("failed to get comments"))
		return
	}

	respond.JSON(c.Writer, http.StatusOK, comments)
}

// Delete deletes the comment with the given ID and all nested descendants.
func (h *Handler) Delete(c *ginext.Context) {
	// Extract id from query params.
	id := c.Param("id")
	if id == "" {
		zlog.Logger.Warn().Msg("id is required")
		respond.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("id is required"))
		return
	}

	// Delete comment.
	err := h.service.DeleteComment(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		// If comment not found, return 404.
		if errors.Is(err, comment.ErrCommentNotFound) {
			zlog.Logger.Error().Err(err).Msg("comment not found")
			respond.Fail(c.Writer, http.StatusNotFound, err)
			return
		}

		zlog.Logger.Error().Err(err).Msg("failed to delete comment")
		respond.Fail(c.Writer, http.StatusInternalServerError, err)
		return
	}

	respond.OK(c.Writer, "comment deleted")
}
