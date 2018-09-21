package routing

import (
	"context"
	"net/http"

	"github.com/billglover/bbot/pkg/queue"
	"github.com/pkg/errors"
)

// Router is responsible for receiving inbound requests, responding to the
// caller and routing themto the appropriate outbound queue.
type Router struct {
	ActionQ   queue.Queuer
	CommandQ  queue.Queuer
	EventQ    queue.Queuer
	ReqSecret string
}

// NewRouter returns a default router.
func NewRouter(actionQueueName, commandQueueName, eventQueueName string) (*Router, error) {
	r := new(Router)

	// Configure one queue for each message type; action, event, and command.
	if actionQueueName != "" {
		q, err := queue.NewSQSQueue(actionQueueName)
		if err != nil {
			return nil, errors.Wrap(err, "ERROR: unable to establish action queue")
		}
		r.ActionQ = q
	}

	if commandQueueName != "" {
		q, err := queue.NewSQSQueue(commandQueueName)
		if err != nil {
			return nil, errors.Wrap(err, "ERROR: unable to establish command queue")
		}
		r.CommandQ = q
	}

	if eventQueueName != "" {
		q, err := queue.NewSQSQueue(eventQueueName)
		if err != nil {
			return nil, errors.Wrap(err, "ERROR: unable to establish event queue")
		}
		r.EventQ = q
	}

	return r, nil
}

// Route handles inbound requests and returns the appropriate response to
// the caller.
func (r *Router) Route(ctx context.Context, req Request) (Response, error) {
	if validateRequest(req, r.ReqSecret) == false {
		return errorResponse("invalid request, check request signature", http.StatusBadRequest)
	}

	switch req.PathParameters["type"] {

	case "event":
		return r.handleEvent(ctx, req)

	case "command":
		return r.handleCommand(ctx, req)

	case "action":
		return r.handleAction(ctx, req)

	default:
		return errorResponse("invalid request, check endpoint type", http.StatusNotFound)
	}
}
