package senders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// SendJson sends jsonData to url using POST and sets the Authorization header to apiToken.
func SendJson(url string, jsonData []byte, apiToken string) error {
	if len(jsonData) == 0 {
		return fmt.Errorf("jsonData is empty")
	}

	maxRetries := 999999999999999999
	backoff := 10 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 100ms, 200ms, 400ms, 800ms, 1600ms
			time.Sleep(backoff * time.Duration(1<<(attempt-1)))
		}

		req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
		if err != nil {
			// Request creation failed, retry
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(jsonData)))
		if apiToken != "" {
			req.Header.Set("Authorization", apiToken)
		}

		response, err := http.DefaultClient.Do(req)
		if err != nil {
			// Connection error (EOF, timeout, etc.), retry
			fmt.Printf("Connection error (attempt %d/%d): %v\n", attempt+1, maxRetries, err)
			continue
		}

		// CRITICAL: Always read and close body to reuse connections
		body, readErr := io.ReadAll(response.Body)
		io.Copy(io.Discard, response.Body) // Drain remaining bytes
		response.Body.Close()

		if readErr != nil {
			// Read error, retry
			fmt.Printf("Read error (attempt %d/%d): %v\n", attempt+1, maxRetries, readErr)
			continue
		}

		// Connection successful, check status code
		if response.StatusCode != http.StatusOK {
			if strings.Contains(string(body), "409") {
				return nil
			}
			// Status code error - panic
			panic(fmt.Sprintf("HTTP error - code: %d message: %s", response.StatusCode, string(body)))
		}

		// Success!
		return nil
	}

	return fmt.Errorf("failed to connect after %d attempts", maxRetries)
}

func RegisterAgent(url string, serverKey string, agentName string) (jwt string, err error) {
	authData := map[string]interface{}{
		"key": serverKey,
		"name": agentName,
	}
	jsonBytes, err := json.Marshal(authData)
	if err != nil {
		panic(err)
	}
	url = fmt.Sprintf("%s/api/v1/agent/management/register", url)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return
	}
	var result string
	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		return
	}
	jwt = result
	return
}