package server

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

// New creates a new HTTP server with the specified address and router.
func New(addr string, router *ginext.Engine) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
