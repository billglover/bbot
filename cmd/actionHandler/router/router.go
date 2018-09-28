/*
Package router provides a service for routing Slack message actions to queues
for processing. It validates all requests using the Slack signing key to ensure
that all requests originated from Slack. Invalid requests are rejected.

The router responds to the original request indicating the message has been
routed successfully (accepted). If it is unable to route the request an
appropriate error response is returned.
*/
package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/billglover/bbot/pkg/agw"
	"github.com/billglover/bbot/pkg/queue"
	"github.com/billglover/bbot/pkg/slack"
)

// Router requires access to the Slack signing secret and the mapping between
// message actions and queues.
type Router struct {
	signingSecret string
	queues        map[string]queue.Queuer
}

// New returns a new Router. It optionally takes configuration functions to
// modify the default configuration.
func New(options ...func(*Router) error) (*Router, error) {
	r := new(Router)
	r.queues = make(map[string]queue.Queuer)
	for _, option := range options {
		err := option(r)
		if err != nil {
			return r, err
		}
	}

	return r, nil
}

// SigningSecret sets the Slack Signing Secret used when validating requests
// during routing.
func SigningSecret(id string) func(*Router) error {
	return func(r *Router) error {
		r.signingSecret = id
		return nil
	}
}

// RegisterRoute associates a mapping between a message action identifier and
// an outbound queue.s
func (r *Router) RegisterRoute(id, url string) error {
	q, err := queue.NewSQSQueue(url)
	if err != nil {
		return err
	}
	r.queues[id] = q
	return nil
}

// Route takes a context and an inbound request. It routes the request to a queue based
// on the registered routes. It returns a response and an error.
func (r *Router) Route(ctx context.Context, req agw.Request) (agw.Response, error) {
	if req.IsValid(r.signingSecret) == false {
		fmt.Println("ERROR: invalid request, check request signature")
		return agw.ErrorResponse("invalid request, check request signature", http.StatusBadRequest)
	}

	action, err := slack.ParseAction(req.Body)
	if err != nil {
		fmt.Println("ERROR: unable to parse message action:", err)
		return agw.ErrorResponse("unable to parse message action", http.StatusBadRequest)
	}

	q, ok := r.queues[action.CallbackID]
	if ok == false {
		fmt.Println("ERROR: message action not supported")
		return agw.ErrorResponse("message action not supported: "+action.CallbackID, http.StatusNotImplemented)
	}

	h := queue.Headers{
		"Team": action.Team.ID,
	}

	err = q.Queue(h, action)
	if err != nil {
		fmt.Println("ERROR: unable to handle message action:", err)
		return agw.ErrorResponse("unable to handle message action", http.StatusInternalServerError)
	}

	fmt.Println("INFO: action queued for processing")
	return agw.SuccessResponse()
}
