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
		name TEXT NOT NULL UNIQUE,
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
	db, err = sql.Open("sqlite3", "/db/api.db")
	if err != nil {
		return
	}
	err = createTables(db)
	return
}

// agent stuff
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

func getAgents(db *sql.DB) (results []map[string]interface{}, err error) {
	query := `SELECT id, name, created_at FROM agents;`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}	
	// Iterate through rows
	for rows.Next() {
		// Create a slice of interface{} to hold each column value
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		
		// Scan the row into the value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		
		// Create a map for this row
		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			
			// Convert []byte to string (common for TEXT fields)
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		
		results = append(results, row)
	}
	
	// Check for errors from iterating
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return results, nil
}

func getAgent(db *sql.DB, agentName string) (result map[string]interface{}, err error) {
	query := `SELECT id, name, created_at FROM agents WHERE name = ?`
	rows, err := db.Query(query, agentName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	// values and value pointers
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	
	for i := range columns {
		valuePtrs[i] = &values[i]
	}
	
	// Scan the row into the value pointers
	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}
	
	// Create a map for this row
	row := make(map[string]interface{})
	for i, col := range columns {
		val := values[i]
		
		// Convert []byte to string (common for TEXT fields)
		if b, ok := val.([]byte); ok {
			row[col] = string(b)
		} else {
			row[col] = val
		}
	}
	result = row
	return
}


func SQL_Delete(db *sql.DB, agentName string) (err error) {
	query := `DELETE FROM agents WHERE name = ?`
	_, err = db.Exec(query, agentName)
	return err
}


// user auth and creation
func CreateUser(db *sql.DB, username, password string) (err error) {
	hashed := sha256.New()
	hashed.Write([]byte(password))
	hashedPass := hex.EncodeToString(hashed.Sum(nil))

	query := `INSERT INTO users (username, hash) VALUES (?, ?);`
	_, err = db.Exec(query, username, hashedPass)
	return err
}

func AuthUser(db *sql.DB, username, password string) (err error) {
	hashed := sha256.New()
	hashed.Write([]byte(password))
	hashedPass := hex.EncodeToString(hashed.Sum(nil))
	var dbHash string
	query := `SELECT hash FROM users WHERE name = ?;`
	err = db.QueryRow(query, username).Scan(&dbHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid password or username!")
		}
		return err
	}
	if dbHash == hashedPass {
		return nil
	}
	return fmt.Errorf("invalid password or username!")
}