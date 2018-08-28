package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type FlagEvent struct {
	
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, evt events.SQSEvent) error {
	for _, msg := range evt.Records {
		fmt.Println("body", msg.Body)
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
