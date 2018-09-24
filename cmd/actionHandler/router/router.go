package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/billglover/bbot/pkg/agw"
	"github.com/billglover/bbot/pkg/queue"
	"github.com/billglover/bbot/pkg/slack"
	"github.com/pkg/errors"
)

// Router handles Message Action requests and validates them before routing
// them to appropriate queues.
type Router struct {
	clientID      string
	clientSecret  string
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

// Route takes a context and an inbound request. It routes the request based on
// the registered routes. It returns a response to the gateway.
func (r *Router) Route(ctx context.Context, req agw.Request) (agw.Response, error) {
	if req.IsValid(r.signingSecret) == false {
		fmt.Println("ERROR: invalid request, check request signature")
		return agw.ErrorResponse("invalid request, check request signature", http.StatusBadRequest)
	}

	action, err := parseAction(req.Body)
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

func parseAction(b string) (slack.MessageAction, error) {
	ma := slack.MessageAction{}

	form, err := url.ParseQuery(b)
	if err != nil {
		return ma, errors.Wrap(err, "failed to parse request body")
	}

	err = json.Unmarshal([]byte(form.Get("payload")), &ma)
	if err != nil {
		return ma, errors.Wrap(err, "failed to parse request body")
	}
	return ma, err
}
