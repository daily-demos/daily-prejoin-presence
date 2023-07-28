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
	"os"
	"time"
)

type Room struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type createParams struct {
	Name       string    `json:"name,omitempty"`
	Properties roomProps `json:"properties,omitempty"`
}

type roomProps struct {
	Exp int64 `json:"exp,omitempty"`
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	apiKey := os.Getenv("DAILY_API_KEY")
	if apiKey == "" {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       util.NewErrorBody("server authentication with Daily failed"),
		}, nil
	}
	room, err := createRoom(apiKey)
	if err != nil {
		errMsg := "failed to create room"
		fmt.Printf("\n%s: %v", errMsg, err)

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       util.NewErrorBody(fmt.Sprintf("%s (check server logs)", errMsg)),
		}, nil
	}
	data, err := json.Marshal(room)
	if err != nil {
		errMsg := "failed to marshal room"
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

// createRoom creates a Daily room.
func createRoom(apiKey string) (*Room, error) {
	name, err := generateNameWithPrefix("presence-")
	if err != nil {
		return nil, fmt.Errorf("failed to generate room name: %w", err)
	}
	params := createParams{
		Name:       name,
		Properties: roomProps{Exp: time.Now().Add(time.Hour).Unix()},
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

// generateNameWithPrefix generates a Daily room name with
// the given prefix.
func generateNameWithPrefix(prefix string) (string, error) {
	s, err := generateRandStr(20)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", prefix, s), nil
}

// generateRandStr generates a string of the request length.
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
