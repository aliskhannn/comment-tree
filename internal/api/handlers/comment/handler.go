package comment

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"

	"github.com/aliskhannn/comment-tree/internal/api/respond"
	"github.com/aliskhannn/comment-tree/internal/model"
	"github.com/aliskhannn/comment-tree/internal/repository/comment"
)

// Service is the interface for the comment service.
type Service interface {
	CreateComment(ctx context.Context, comment *model.Comment) (model.Comment, error)
	GetCommentsByParentID(ctx context.Context, parentID uuid.UUID) ([]model.Comment, error)
	GetComments(ctx context.Context, parentID *uuid.UUID, search, sort string, limit, offset int) ([]model.Comment, error)
	DeleteComment(ctx context.Context, id uuid.UUID) error
}

// Handler is the handler for the comment API.
type Handler struct {
	service Service
}

// NewHandler creates a new Handler.
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateRequest is the request for the create comment API.
type CreateRequest struct {
	ParentID *uuid.UUID `json:"parent_id"`
	Content  string     `json:"content" binding:"required,min=1,max=1000"`
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

	zlog.Logger.Printf("parent id before: %v", req.ParentID)

	var parentID *uuid.UUID
	if req.ParentID == nil {
		parentID = nil
	} else {
		parentID = req.ParentID
	}

	// Create the comment.
	cm := &model.Comment{
		ParentID: parentID,
		Content:  req.Content,
	}

	zlog.Logger.Printf("parent id in comment after: %v", parentID)

	res, err := h.service.CreateComment(c, cm)
	if err != nil {
		respond.Fail(c.Writer, http.StatusInternalServerError, err)
		return
	}

	respond.Created(c.Writer, res)
}

// GetTree retrieves the comment with the given ID and all nested descendants.
func (h *Handler) GetTree(c *ginext.Context) {
	// Extract id from query params.
	idStr := c.Param("id")
	if idStr == "" {
		zlog.Logger.Warn().Msg("parent id is required")
		respond.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("parent id is required"))
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to parse parent id")
		respond.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid parent id"))
		return
	}

	// Get comments.
	comments, err := h.service.GetCommentsByParentID(c.Request.Context(), id)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to get comments")
		respond.Fail(c.Writer, http.StatusInternalServerError, fmt.Errorf("failed to get comments"))
		return
	}

	respond.JSON(c.Writer, http.StatusOK, comments)
}

// GetList retrieves comments with pagination, sorting, and optional search.
func (h *Handler) GetList(c *ginext.Context) {
	// Get query params.
	parentIDStr := c.Query("parent")
	var parentID *uuid.UUID
	if parentIDStr != "" {
		id, err := uuid.Parse(parentIDStr)
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("failed to parse parent id")
			respond.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid parent id"))
			return
		}
		parentID = &id
	}

	search := c.Query("search")
	sort := c.DefaultQuery("sort", "created_at_asc")

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	comments, err := h.service.GetComments(c.Request.Context(), parentID, search, sort, limit, offset)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to get comments")
		respond.Fail(c.Writer, http.StatusInternalServerError, fmt.Errorf("failed to get comments"))
		return
	}

	respond.JSON(c.Writer, http.StatusOK, comments)
}

// Delete deletes the comment with the given ID and all nested descendants.
func (h *Handler) Delete(c *ginext.Context) {
	// Extract id from query params.
	idStr := c.Param("id")
	if idStr == "" {
		zlog.Logger.Warn().Msg("parent id is required")
		respond.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("parent id is required"))
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to parse parent id")
		respond.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid parent id"))
		return
	}

	// Delete comment.
	err = h.service.DeleteComment(c.Request.Context(), id)
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
