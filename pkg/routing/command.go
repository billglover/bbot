package routing

import (
	"context"
	"net/http"
)

func (r *Router) handleCommand(ctx context.Context, req Request) (Response, error) {
	return errorResponse("command API not yet implemented", http.StatusNotImplemented)
}
