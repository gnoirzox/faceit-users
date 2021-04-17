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

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func OpenDBConnection() *sql.DB {
	var (
		host     = getEnv("DB_HOST", "localhost")
		port     = getEnv("DB_PORT", "5438")
		user     = getEnv("DB_USER", "postgres")
		password = getEnv("DB_PASS", "postgres")
		dbname   = getEnv("DB_NAME", "postgres")
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Println(err.Error())
	}

	return db
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
