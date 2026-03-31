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

	//	API_AUTH "github.com/Ceald1/purrimeter/api/routes/auth"
	//	API_LOGGING "github.com/Ceald1/purrimeter/api/routes/logging"
	API_MANAGEMENT "github.com/Ceald1/purrimeter/api/routes/management"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
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

func LOAD() {
	SURREAL_ADMIN = os.Getenv("SURREAL_ADMIN")
	SURREAL_PASS = os.Getenv("SURREAL_PASS")
	API_ADMIN_USER = os.Getenv("API_ADMIN")
	API_ADMIN_PASS = os.Getenv("API_ADMIN_PASS")
	AGENT_SECRET = os.Getenv("AGENT_SECRET")
	SURREAL_DB_HOST = os.Getenv("SURREAL_HOST")
	JWT_SECRET = os.Getenv("ENRICHMENT_JWT")
	FAKE_AGENT_ID = os.Getenv("FAKE_AGENT_ID")
}

func main() {
	godotenv.Load()
	LOAD()
	var err error
	var result_data map[string]interface{}
	log.Info("starting tests....")

	log.Info("beginning API tests....")

	log.Info("starting management tests.....")
	log.Info(fmt.Sprintf("logging in with %s:%s user...", SURREAL_ADMIN, SURREAL_PASS))
	root_url := "http://127.0.0.1:8080"
	login_url := fmt.Sprintf("%s/api/v2/management/user/login", root_url)
	loginData := API_MANAGEMENT.UserLogin{}
	loginData.Username = SURREAL_ADMIN
	loginData.Password = SURREAL_PASS
	err, result_data = MakePOSTRequest(login_url, loginData, "")
	if err != nil {
		log.Fatal(err)
	}
	userToken := result_data["result"]
	log.Info(userToken)

}

func MakePOSTRequest(url string, bodyData any, apiToken string) (err error, result_data map[string]interface{}) {
	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		return
	}
	requestBody := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return
	}
	if apiToken != "" {
		req.Header.Set("Authorization", apiToken)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if strings.Contains(string(body), "error") {
		return errors.New(string(body)), nil
	}
	err = json.Unmarshal(body, &result_data)

	return
}
