package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/daily-demos/daily-prejoin-presence/m/v2/util"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Room struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type createParams struct {
	Name string `json:"name"`
	Exp  int64  `json:"exp,omitempty"`
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	q, err := url.ParseQuery(request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}
	roomName := q.Get("roomName")
	if roomName == "" {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "roomName parameter not found in request",
		}, nil
	}

	apiKey := os.Getenv("DAILY_API_KEY")
	if apiKey == "" {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "server authentication with Daily failed",
		}, nil
	}

	room, err := createRoom(apiKey)
	if err != nil {
		errMsg := "failed to marshal participants"
		fmt.Printf("\n%s: %v", errMsg, err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("%s (check server logs)", errMsg),
		}, nil
	}
	data, err := json.Marshal(room)
	if err != nil {
		errMsg := "failed to marshal participants"
		fmt.Printf("\n%s: %v", errMsg, err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("%s (check server logs)", errMsg),
		}, nil
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(data),
	}, nil
}

// CreateWithPrefix creates a room with the name containing the specified
// prefix. The rest of the name is randomized.
func createRoom(apiKey string) (*Room, error) {
	name, err := generateNameWithPrefix("presence-")
	if err != nil {
		return nil, fmt.Errorf("failed to generate room name: %w", err)
	}
	params := createParams{
		Name: name,
		Exp:  time.Now().Add(time.Hour).Unix(),
	}
	params.Name = name

	// Make the request body for room creation
	reqBody, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to make room creation request body: %w", err)
	}

	endpoint := fmt.Sprintf("%s/rooms", util.DailyAPIURL)

	// Make the actual HTTP request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request to rooms endpoint: %w", err)
	}

	// Prepare auth and content-type headers for request
	util.SetAPIKeyAuthHeaders(req, apiKey)

	// Do the thing!!!
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}

	// Parse the response
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read room creation response body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to make API call to Daily: %d: %s", res.StatusCode, string(resBody))
	}

	var room Room
	if err := json.Unmarshal(resBody, &room); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Daily response to room: %w", err)
	}

	return &room, nil
}

func generateNameWithPrefix(prefix string) (string, error) {
	s, err := generateRandStr(20)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", prefix, s), nil
}

func generateRandStr(length int) (string, error) {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[num.Int64()]
	}
	return string(result), nil
}
