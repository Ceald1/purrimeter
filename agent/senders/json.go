package senders

import (
	"bytes"
	"net/http"
	"fmt"
	"io/ioutil"
)


func SendJson(url string, jsonData []byte) (err error) {
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		err = fmt.Errorf("code: %d message: %s", response.StatusCode, string(body))
	}
	return err
}