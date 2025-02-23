package driver

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	//_ "github.com/jackc/v5/pgconn"
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifetime)

	dbConn.SQL = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// testDB tries to ping the Database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}

	return err

}

// NewDatabase creates a new database connection for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

type OllamaDriver struct {
	Host   string
	Model  string
	Stream string
}

func NewOllamaDriverStream() *OllamaDriver {
	return &OllamaDriver{
		Host:   os.Getenv("LLM_HOST"),
		Model:  os.Getenv("LLM_MODEL"),
		Stream: "true",
	}
}

func NewOllamaDriverNoStream() *OllamaDriver {
	fmt.Println("Running NewOllamaDriverNoStream")
	return &OllamaDriver{
		Host:   os.Getenv("LLM_HOST"),
		Model:  os.Getenv("LLM_MODEL"),
		Stream: "false",
	}
}

// func (o *OllamaDriver) OllamaGenerateResponse(prompt string) (string, error) {

// 	fmt.Println("OllamaGenerateResponse")
// 	request := map[string]interface{}{
// 		"model":  o.Model,
// 		"prompt": prompt,
// 		"stream": o.Stream,
// 	}

// 	jsonData, _ := json.Marshal(request)

// 	resp, err := http.Post(o.Host+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		fmt.Println("error calling LLM: ", err)
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("error reading response body: %w", err)
// 	}

// 	// Print raw response for debugging
// 	fmt.Println("Raw Ollama Response:", string(body))

// 	var response map[string]interface{}

// 	// Parse JSON response
// 	if err := json.Unmarshal(body, &response); err != nil {
// 		return "", fmt.Errorf("error decoding JSON: %w", err)
// 	}

// 	// Print the parsed response (for debugging)
// 	fmt.Printf("Parsed JSON Response: %+v\n", response)

// 	respText, ok := response["response"].(string)
// 	if !ok {
// 		return "", fmt.Errorf("unexpected response format: %v", response)
// 	}

// 	return respText, nil
// }
