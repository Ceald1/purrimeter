package senders

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func SendJson(url string, jsonData []byte) (err error) {
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		err = fmt.Errorf("code: %d message: %s", response.StatusCode, string(body))
	}
	return err
}
