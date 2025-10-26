package db

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/valkey-io/valkey-go"
)
var ctx = context.Background()

func Valkey_Init() (client valkey.Client, err error) {
	pass := os.Getenv("REDIS_PASS")
	client, err = valkey.NewClient(valkey.ClientOption{InitAddress: []string{"valkey:6379"}, Password: pass})
	return
}


func Valkey_Secrets(client valkey.Client) (err error) {
	client.Do(ctx, client.B().Flushall().Build())
	var secret_key string
	db, err := SQL_Init()
	if err != nil {
		return err
	}
	secret_key, err = getSecret(db)
	if err != nil {
		return err
	}

	err = client.Do(ctx, client.B().Set().Key("agentKey").Value(secret_key).Nx().Build()).Error()
	return
}

func Get_Valkey_Secrets(client valkey.Client) (secret_key string, err error){
	secret_key, err = client.Do(ctx, client.B().Get().Key("agentKey").Build()).ToString()
	return
}

func Valkey_MoveAgentsToDB(client valkey.Client) (err error) {
	db, err := SQL_Init()
	if err != nil {
		return err
	}
	agents, err := getAgents(db)
	if err != nil {
		return err
	}
	for _, agent := range agents {
		key := fmt.Sprintf("agent:%s", agent["name"])
		fields := make(map[string]string)
		for k, v := range agent {
			fields[k] = fmt.Sprintf("%v", v)
		}

		cmd := client.B().Hset().Key(key).FieldValue()
		for field, value := range fields {
			cmd = cmd.FieldValue(field, value)
		}
		err = client.Do(ctx, cmd.Build()).Error()
		if err != nil {
			return err
		}
	}
	return
}



func ValkeyAgentToDB(client valkey.Client, agentName string) (err error) {
	db, err := SQL_Init()
	if err != nil {
		return err
	}
	agent, err := getAgent(db, agentName)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("agent:%s", agent["name"])
	fields := make(map[string]string)
	for k, v := range agent {
		fields[k] = fmt.Sprintf("%v", v)
	}
	cmd := client.B().Hset().Key(key).FieldValue()
	for field, value := range fields {
		cmd = cmd.FieldValue(field, value)
	}
	err = client.Do(ctx, cmd.Build()).Error()
	return err

}

func Valkey_FetchAgent(client valkey.Client, agentName string) (err error) {
	key := fmt.Sprintf("agent:%s", strings.ToLower(agentName))
	result, err := client.Do(ctx, client.B().Exists().Key(key).Build()).AsInt64()
	if err != nil {
		return err
	}
	if result > 0 {
		return nil
	}
	return fmt.Errorf("agent doesn't exist!")
}