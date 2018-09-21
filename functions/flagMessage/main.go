package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/billglover/bbot/pkg/flagging"
)

func main() {
	lambda.Start(flagging.Handler)
}
