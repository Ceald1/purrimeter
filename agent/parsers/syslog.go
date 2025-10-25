package parsers

// for rsyslog

import (
	"fmt"

	"encoding/json"
	"regexp"
	"strings"
	"time"
	"github.com/matishsiao/goInfo"

	"github.com/Ceald1/purrimeter/agent/senders"
	"github.com/nxadm/tail"
	rus "github.com/sirupsen/logrus"
)
var syslogRegex = regexp.MustCompile(`^(\S+)\s+(\S+)\s+([^:\[\s]+)(?:\[(\d+)\])?:\s*(.*)$`)

type SyslogMessage struct {
	Timestamp time.Time
	Host      string
	Program   string
	PID       string
	Message   string
}

func parseSyslogLine(line string) (*SyslogMessage, error) {
	matches := syslogRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, nil // Not a valid syslog line
	}

	// Parse timestamp (RFC3339 format with timezone)
	timestamp, err := time.Parse(time.RFC3339Nano, matches[1])
	if err != nil {
		timestamp = time.Now() // Fallback to current time
	}

	return &SyslogMessage{
		Timestamp: timestamp,
		Host:      matches[2],
		Program:   matches[3],
		PID:       matches[4], // May be empty
		Message:   strings.TrimSpace(matches[5]),
	}, nil
}

func RsyslogTail(logPath string, url string) (err error) {
	// tails log file, decodes to json, and sends line to server
	
	t, err := tail.TailFile(logPath, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		return
	}
	logger := rus.New()
	for l := range t.Lines {
		line := l.Text
		if strings.TrimSpace(line) == "" {
			continue
		}
		json_log := decodeLog(line, logger)
		if json_log["timestamp"] != nil {
			jsonData, err := json.Marshal(json_log)
			if err != nil {
				break
			}
			if json_log["program"] != "sysmon"{
				err = senders.SendJson(url, jsonData)
				if err != nil {
					fmt.Println(err.Error())
				}
			}else{
				re_html := regexp.MustCompile(`&#x[A-Fa-f0-9]{0,2};`)
				syslog := json_log["message"].(string)
				syslog = strings.ToValidUTF8(syslog, "")
				syslog = re_html.ReplaceAllString(syslog, "_")
				sysmon_log := XML(syslog)
				new_json_log := make(map[string]interface{})
				new_json_log["timestamp"] = json_log["timestamp"]
				new_json_log["program"]   = json_log["program"]
				new_json_log["message"]   = sysmon_log
				new_json_log["pid"]       = json_log["pid"]
				new_json_log["host"]	  = json_log["host"]
				new_json_log["logOrigin"] = "sysmon"
				new_json_log["os"]  = json_log["os"]
				jsonData, _= json.Marshal(new_json_log)
				err = senders.SendJson(url, jsonData)
				if err != nil {
					fmt.Println(err.Error())
				}



			}

		}

	}
	return
}


func decodeLog(log string, logger *rus.Logger)(output map[string]interface{}){
	// decode log to prep it to be sent as json
	info, _ := goInfo.GetInfo()
	log_message, err := parseSyslogLine(log)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	logger.SetFormatter(&rus.JSONFormatter{
		TimestampFormat:time.RFC3339,
	})
	logger.SetLevel(rus.InfoLevel)
	entry := rus.Fields{
		"timestamp": 	log_message.Timestamp,
		"program":   	log_message.Program,
		"message":   	log_message.Message,
		"pid":       	log_message.PID,
		"host":      	info.Hostname,
		"logOrigin":    "syslog",
		"os":     info.GoOS,
	}
	output = map[string]interface{}(entry)
	return
}