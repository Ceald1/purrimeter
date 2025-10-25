package db

import (
	"github.com/valkey-io/valkey-go"
	"os"
	"context"
)
var ctx = context.Background()

func Valkey_Init() (client valkey.Client, err error) {
	pass := os.Getenv("REDIS_PASS")
	client, err = valkey.NewClient(valkey.ClientOption{InitAddress: []string{"valkey:6379"}, Password: pass})
	return
}


func Valkey_Secrets(client valkey.Client) (err error) {
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