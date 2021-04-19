package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
)

func GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func OpenDBConnection() (*sql.DB, error) {
	var (
		host     = GetEnv("DB_HOST", "localhost")
		port     = GetEnv("DB_PORT", "5438")
		user     = GetEnv("DB_USER", "postgres")
		password = GetEnv("DB_PASS", "postgres")
		dbname   = GetEnv("DB_NAME", "postgres")
	)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Println("%s: %s", "Failed to connect to Postgres", err)

		return nil, err
	}

	return db, nil
}

func ReturnJsonResponse(w http.ResponseWriter, httpCode int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		log.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(response)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
