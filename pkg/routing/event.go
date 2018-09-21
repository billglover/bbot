package routing

import (
	"context"
	"net/http"
)

func (r *Router) handleEvent(ctx context.Context, req Request) (Response, error) {
	return errorResponse("events API not yet implemented", http.StatusNotImplemented)
}
