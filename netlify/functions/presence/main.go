package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/daily-demos/daily-prejoin-presence/m/v2/util"
	"io"
	"net/http"
	"os"
)

// Participant represents a Daily participant.
// who is already in a Daily room
type Participant struct {
	ID   string `json:"id"`
	Name string `json:"userName"`
}

// presenceRes is the expected response format
// from Daily's REST API when retrieving room presence.
type presenceRes struct {
	Data []Participant
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	q := request.QueryStringParameters
	roomName, ok := q["roomName"]
	if !ok {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       util.NewErrorBody("roomName parameter not found in request"),
		}, nil
	}

	apiKey := os.Getenv("DAILY_API_KEY")
	if apiKey == "" {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       util.NewErrorBody("server authentication with Daily failed"),
		}, nil
	}

	participants, err := getPresence(roomName, apiKey)
	if err != nil {
		errMsg := "failed to get presence"
		fmt.Printf("\n%s: %v", errMsg, err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       util.NewErrorBody(fmt.Sprintf("%s (check server logs)", errMsg)),
		}, nil
	}
	data, err := json.Marshal(participants)
	if err != nil {
		errMsg := "failed to marshal participants"
		fmt.Printf("\n%s: %v", errMsg, err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       util.NewErrorBody(fmt.Sprintf("%s (check server logs)", errMsg)),
		}, nil
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(data),
	}, nil
}

// getPresence retrieves all participants already in the given Daily room
func getPresence(roomName string, apiKey string) ([]Participant, error) {
	endpoint := fmt.Sprintf("%s/rooms/%s/presence", util.DailyAPIURL, roomName)
	// Make the actual HTTP request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request to room endpoint: %w", err)
	}

	util.SetAPIKeyAuthHeaders(req, apiKey)

	// Do the thing!!!
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	// Parse the response
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read presence response body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to make API call to Daily: %d: %s", res.StatusCode, string(resBody))
	}

	var pr presenceRes
	if err := json.Unmarshal(resBody, &pr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Daily response to participant slice: %w", err)
	}
	return pr.Data, nil
}
