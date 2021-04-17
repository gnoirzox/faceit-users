package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"./users"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/user", users.CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}", users.UpdateUser).Methods("PUT")
	router.HandleFunc("/user/{id}", users.DeleteUser).Methods("DELETE")
	router.HandleFunc("/user/{id}", users.GetUser).Methods("GET")
	router.HandleFunc("/users", users.GetUsers).Methods("GET")

	log.Fatal(http.ListenAndServe(":8888", router))
}
