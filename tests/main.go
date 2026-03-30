package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	API_AUTH "github.com/Ceald1/purrimeter/api/routes/auth"
	API_LOGGING "github.com/Ceald1/purrimeter/api/routes/logging"
	API_MANAGEMENT "github.com/Ceald1/purrimeter/api/routes/management"
	"github.com/charmbracelet/log"
)

var (
	SURREAL_ADMIN   string = os.Getenv("SURREAL_ADMIN")
	SURREAL_PASS    string = os.Getenv("SURREAL_PASS")
	API_ADMIN_USER  string = os.Getenv("API_ADMIN")
	API_ADMIN_PASS  string = os.Getenv("API_ADMIN_PASS")
	AGENT_SECRET    string = os.Getenv("AGENT_SECRET")
	SURREAL_DB_HOST string = os.Getenv("SURREAL_HOST")
	JWT_SECRET      string = os.Getenv("ENRICHMENT_JWT")
	FAKE_AGENT_ID   string = os.Getenv("FAKE_AGENT_ID")
)

func main() {
	log.Info("starting tests....")

	log.Info("beginning API tests....")

}

func MakePOSTRequest(url string, bodyData any, apiToken string) (err error) {
	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		return err
	}
	requestBody := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", apiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if strings.Contains(string(body), "error") {
		return errors.New(string(body))
	}

	return nil
}
