package queue

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSEvent is a Queue event.
type SQSEvent events.SQSEvent

// SQSQueue implements the Queue interface.
type SQSQueue struct {
	svc  *sqs.SQS
	name string
}

// Queue takes message headers and a body and places it onto the SQS queue.
func (q *SQSQueue) Queue(h Headers, b Body) error {

	body, err := json.Marshal(b)
	if err != nil {
		return err
	}

	delay := aws.Int64(0)
	attributes := make(map[string]*sqs.MessageAttributeValue)

	for k, v := range h {
		attributes[k] = &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(v),
		}
	}

	msg := &sqs.SendMessageInput{
		DelaySeconds:      delay,
		MessageAttributes: attributes,
		MessageBody:       aws.String(string(body)),
		QueueUrl:          aws.String(q.name),
	}

	_, err = q.svc.SendMessage(msg)
	return err
}

// NewSQSQueue takes the name of an AWS SQS queue and returns a pointer to a Q.
func NewSQSQueue(name string) (*SQSQueue, error) {
	if name == "" {
		return nil, errors.New("Queue name cannot be empty")
	}
	q := new(SQSQueue)
	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	q.svc = sqs.New(sess)
	q.name = name
	return q, nil
}
