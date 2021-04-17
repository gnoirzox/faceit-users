package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"../utils"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var userId string

	if Id, ok := vars["id"]; ok {
		userId = string(Id)
	} else {
		log.Println(err.Error())

		utils.ReturnJsonResponse(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "The provided user id is not valid."}
		)

		return
	}

	user, err := RetrieveUser(userId)

	if err != nil {
		log.Println(err.Error())

		utils.ReturnJsonResponse(
			w,
			http.StatusNotFound,
			map[string]string{"error": "The user does not exist."}
		)

		return
	}

	utils.ReturnJsonResponse(w, http.StatusOK, user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filters :=  make(map[string]string)

	if nickname, ok := vars["nickname"]; ok {
		filters["nickname"] = nickname

	}

	if email, ok := vars["email"]; ok {
		filters["email"] = email

	}

	if country, ok := vars["country"]; ok {
		filters["country"] = country

	}

	users, err := RetrieveUsers(filters)

	if err != nil {
		log.Println(err.Error())

		utils.ReturnJsonResponse(
			w,
			http.StatusNotFound,
			map[string]string{"error": "The given filters don't seem valid."}
		)

		return
	}

	utils.ReturnJsonResponse(w, http.StatusOK, users)
}

func ValidateUser(user *User) bool {
	if !&u.IsValidNickname() {
		utils.ReturnJsonResponse(
			w, 
			http.StatusBadRequest,
			map[string]string{
				"error": "The submitted nickname is not valid. It should have a length between 3 and 12 characters."
			}
		)

		return false
	}

	if !&u.IsValidPassword() {
		utils.ReturnJsonResponse(
			w, 
			http.StatusBadRequest,
			map[string]string{
				"error": "The submitted password is not valid. It should have a length of at least 8 characters."
			}
		)

		return false
	}

	if !&u.IsValidEmail() {
		utils.ReturnJsonResponse(
			w, 
			http.StatusBadRequest,
			map[string]string{"error": "The submitted Email is not valid."}
		)

		return false
	}

	return true
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	var u User

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)

	defer r.Body.Close()

	if err != nil {
		log.Println(err.Error())

		utils.ReturnJsonResponse(
			w, 
			http.StatusBadRequest, 
			map[string]string{"error": "Invalid resquest payload"}
		)

		return
	}

	if !ValidateUser(&u) {
		return
	}

	err = InsertUser(&u)

	if err != nil {
		log.Println(err.Error())

		utils.ReturnJsonResponse(
			w, 
			http.StatusInternalServerError, 
			map[string]string{"status": "Unexpected error"}
		)

		return
	}

	utils.ReturnJsonResponse(
		w, 
		http.StatusOK, 
		map[string]string{"status": "The user has successfully been created"}
	)
}

func PutUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var userId string

	if Id, ok := vars["id"]; ok {
		userId = string(Id)
	} else {
		log.Println(err.Error())

		utils.ReturnJsonResponse(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "The provided user id is not valid."}
		)

		return
	}

	var user User

	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&user)

	defer r.Body.Close()

	if err != nil {
		log.Println(err.Error())

		utils.ReturnJsonResponse(
			w, 
			http.StatusBadRequest, 
			map[string]string{"error": "Invalid resquest payload"}
		)

		return
	}

	if !ValidateUser(&user) {
		return
	}

	err = UpdateUser(&user)

	if err != nil {
		log.Println(err.Error())

		utils.ReturnJsonResponse(
			w, 
			http.StatusInternalServerError, 
			map[string]string{"status": "Unexpected error"}
		)

		return
	}

	utils.ReturnJsonResponse(
		w, 
		http.StatusOK, 
		map[string]string{"status": "The user has successfully been updated"}
	)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var userId string

	if Id, ok := vars["id"]; ok {
		userId = string(Id)
	} else {
		log.Println(err.Error())

		utils.ReturnJsonResponse(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "The provided user id is not valid or not given."}
		)

		return
	}

	user, err := RemoveUser(userId)

	if err != nil {
		log.Println(err.Error())

		utils.ReturnJsonResponse(
			w,
			http.StatusNotFound,
			map[string]string{"error": "The user does not exist."}
		)

		return
	}

	utils.ReturnJsonResponse(
		w,
		http.StatusOK,
		map[string]string{"status": "The provided user has been successfully deleted."}
	)
}
