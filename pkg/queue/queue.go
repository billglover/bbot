package queue

import "context"

// Queuer represents an infrastructure Queue on which messages can be sent. It
// abstracts away the specifics of the specific infrastructure provider.
type Queuer interface {
	Queue(ctx context.Context, h Headers, b Body) error
}

// Headers is a map that contains key value pairs representing message headers.
type Headers map[string]string

// Body is an interface type representing the message body. It is marshalled to/from
// JSON during enqueue/dequeue operaions.
type Body interface{}
