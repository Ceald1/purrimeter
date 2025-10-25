package db

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"

	"crypto/rand"
	"crypto/sha256"

	_ "github.com/mattn/go-sqlite3"
)


type AgentRegisterDB struct {
	ID	int 
	Name string
}

func createTables(db *sql.DB) (err error) {
	// create tables if they don't exist and set default values
	query := `
		CREATE TABLE IF NOT EXISTS agents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(query)
	if err != nil {
		return
	}
	query = `
		CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(query)
	if err != nil {
		return
	}


	secret_key_bytes := make([]byte, 64)
	rand.Read(secret_key_bytes)
	hash := sha256.Sum256(secret_key_bytes)
	secret_key := hex.EncodeToString(hash[:])

	// Create table
	query = `
		CREATE TABLE IF NOT EXISTS secrets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			agentKey TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`
	_, err = db.Exec(query)
	if err != nil {
		return
	}

	// Insert only if no secrets exist
	insertQuery := `
		INSERT INTO secrets (agentKey) 
		SELECT ? 
		WHERE NOT EXISTS (SELECT 1 FROM secrets LIMIT 1)`
	_, err = db.Exec(insertQuery, secret_key)
	return
}




func SQL_Init() (db *sql.DB, err error){
	db, err = sql.Open("sqlite3", "./api.db")
	if err != nil {
		return
	}
	err = createTables(db)
	return
}


func RegisterAgent(db *sql.DB, agentName string) error {
	name := strings.ToLower(agentName)
	
	// Try to insert directly
	query := `INSERT INTO agents (name) VALUES (?);`
	_, err := db.Exec(query, name)
	
	if err != nil {
		// Check if it's a unique constraint violation
		if strings.Contains(err.Error(), "UNIQUE constraint failed") { 
			return fmt.Errorf("agent already exists")
		}
		return err  // Other database error
	}
	
	return nil
}

func getSecret(db *sql.DB) (secretKey string, err error) {
	query := `SELECT agentKey FROM secrets ORDER by id ASC LIMIT 1`
	err = db.QueryRow(query).Scan(&secretKey)

	return
}
