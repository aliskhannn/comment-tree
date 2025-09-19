package router

import (
	"github.com/wb-go/wbf/ginext"

	"github.com/aliskhannn/comment-tree/internal/api/handlers/comment"
	"github.com/aliskhannn/comment-tree/internal/middleware"
)

// New creates a new Gin engine with routes and middlewares for the comment API.
func New(handler *comment.Handler) *ginext.Engine {
	e := ginext.New()

	e.Use(middleware.CORSMiddleware())
	e.Use(ginext.Logger())
	e.Use(ginext.Recovery())

	{
		api := e.Group("/api/comments")
		api.POST("/", handler.Create)
		api.GET("/:id", handler.Get)
		api.GET("/", handler.GetList) // with query params ?parent=&search=&limit=&offset
		api.DELETE("/:id", handler.Delete)
	}

	return e
}
