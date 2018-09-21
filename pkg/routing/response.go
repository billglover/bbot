package routing

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

// Response is of type APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

// ErrorResponse takes a string and generates an appropriate Response.
func errorResponse(msg string, status int) (Response, error) {
	p := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  strconv.Itoa(status),
		Message: msg,
	}

	body, _ := json.Marshal(p)
	resp := Response{
		Body:       string(body),
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
	return resp, nil
}

// SuccessResponse takes a string and generates an appropriate Response.
func successResponse(msg string) (Response, error) {
	p := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  strconv.Itoa(http.StatusAccepted),
		Message: msg,
	}

	body, _ := json.Marshal(p)
	resp := Response{
		Body:       string(body),
		StatusCode: http.StatusAccepted,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
	return resp, nil
}
