package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

var (
	// ErrFailedDailyAPICall is returned when an API call to Daily failed.
	ErrFailedDailyAPICall = errors.New("failed API call to Daily")
)

type errorBody struct {
	Error string `json:"error"`
}

// NewErrorBody returns a JSON string with
// the specified error message.
func NewErrorBody(msg string) string {
	body, _ := json.Marshal(errorBody{
		Error: msg,
	})
	return string(body)
}

// NewNoAPIKeyRes returns an error response when the server
// could not retrieve a Daily API key.
func NewNoAPIKeyRes() *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       NewErrorBody("server authentication with Daily failed"),
	}
}

// NewErrFailedDailyAPICall takes the given error
// and wraps it with a standard API call error
func NewErrFailedDailyAPICall(err error) error {
	return fmt.Errorf("%s: %w", err, ErrFailedDailyAPICall)
}
