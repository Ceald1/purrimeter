package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"encoding/json"
	"net/http"
	"time"

	Framework "github.com/Ceald1/purrimeter/modules/framework"
)

// orchestra microservice/module this module is for orchestrating where alerts go. For registering modules/integrations check the conductor component.

func main(){
	conductor_host := os.Getenv("CONDUCTOR_HOST")
	if conductor_host == "" {
		conductor_host = "127.0.0.1"
	}
	port := os.Getenv("CONDUCTOR_PORT")
	if port == "" {
		port = "8000"
	}
	manticoreHost := os.Getenv("MANTICORE_HOST")
	if manticoreHost == "" {
		manticoreHost = "127.0.0.1"
	}
	manticorePort := os.Getenv("MANTICORE_PORT")
	if manticorePort == "" {
		manticorePort = "9308"
	}
	var Manticore = Framework.NewDBModulehandler(fmt.Sprintf("http://%s:%s", manticoreHost, manticorePort))
	var baseURL = fmt.Sprintf(`http://%s:%s`, conductor_host, port)
	_, err := http.NewRequest(`GET`, fmt.Sprintf(`%s/sync`, baseURL), nil)
	if err != nil {
		panic(err)
	}
	NonEnrichment, err := GetNonEnrichment(baseURL)
	if err != nil {
		panic(err)
	}
	Enrichment, err := GetEnrichment(baseURL)
	if err != nil {
		panic(err)
	}
	// get last query
	lastQueryFile, err := os.OpenFile("query.sql", os.O_RDWR|os.O_CREATE, 0644)
	var query string
	if err != nil {
		panic(err)
	}
	defer lastQueryFile.Close()

	// Read existing query or use default
	stat, _ := lastQueryFile.Stat()
	if stat.Size() == 0 {
		fmt.Println("no last query found....")
		query = `SELECT * FROM purrimeter_alerts LIMIT 100 OFFSET 0;`
		lastQueryFile.WriteString(query)
	} else {
		buf := make([]byte, stat.Size())
		n, err := lastQueryFile.Read(buf)
		if err != nil {
			panic(err)
		}
		query = strings.TrimSpace(string(buf[:n]))
	}
	fmt.Println(query)
	for {
		queryResults, err := Manticore.SQL(query)
		Update(baseURL)
		if err != nil {
			fmt.Println(err.Error())
			if strings.Contains(err.Error(), "EOF") {
				fmt.Println("end of data, eepy time!")
				time.Sleep(5000 * time.Millisecond)
				continue
			}
		}
		if queryResults == nil {
			time.Sleep(5000 * time.Millisecond)
			continue
		}
		logs_raw := queryResults["hits"].(map[string]interface{})["hits"].([]interface{})

		// Convert []interface{} to []map[string]interface{}
		logs := make([]map[string]interface{}, 0, len(logs_raw))
		for _, logInterface := range logs_raw {
			if logMap, ok := logInterface.(map[string]interface{}); ok {
				logs = append(logs, logMap)
			}
		}
		if len(logs) < 1 {
			time.Sleep(5000 * time.Millisecond)
			continue
		}
		for _, log := range logs {
			log_number := log["_source"].(map[string]interface{})["alert_num"].(float64)
			jsonBytes, _ := json.Marshal(log)
			for _, enrichment := range Enrichment {
				enrichmentUrl := enrichment["Module"]["Url"]

				http.Post(enrichmentUrl, "application/json", bytes.NewBuffer(jsonBytes))
				
			}
			for m := range NonEnrichment {
				for _, sub_m := range NonEnrichment[m] {
					mod_url := sub_m["Module"]["Url"]

					http.Post(mod_url, "application/json", bytes.NewBuffer(jsonBytes))

				}
			}
			query = fmt.Sprintf(`SELECT * FROM purrimeter_alerts WHERE log_number > %d LIMIT 200;`, int64(log_number))
			// Truncate and write new query
			lastQueryFile.Truncate(0)
			lastQueryFile.Seek(0, 0)
			lastQueryFile.WriteString(query)
			
		}
		time.Sleep(800 * time.Millisecond)


	}

}


func GetEnrichment(baseurl string) (data []map[string]map[string]string, err error) {
	url := fmt.Sprintf("%s/module/enrichment", baseurl)
	response, err := http.Get(url)
	if err != nil {
		return
	}
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return
	}
	return 
}

func GetNonEnrichment(baseurl string) (data map[string][]map[string]map[string]string, err error) {
	data = make(map[string][]map[string]map[string]string)
	url := fmt.Sprintf("%s/module/!enrichment", baseurl)
	response, err := http.Get(url)
	if err != nil {
		return
	}
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return
	}
	return 
}

func Update(baseurl string) (err error) {
	url := fmt.Sprintf("%s/sync", baseurl)
	_, err = http.Get(url)
	return

}