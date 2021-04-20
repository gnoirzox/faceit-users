package users

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gnoirzox/faceit-users/events"
	"github.com/gnoirzox/faceit-users/utils"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var userId string

	if Id, ok := vars["id"]; ok {
		userId = string(Id)
	} else {
		utils.ReturnJsonResponse(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "The provided user id is not valid."},
		)

		return
	}

	user, err := RetrieveUser(userId)

	if err != nil {
		log.Println(err.Error())

		utils.ReturnJsonResponse(
			w,
			http.StatusNotFound,
			map[string]string{"error": "The user does not exist."},
		)

		return
	}

	utils.ReturnJsonResponse(w, http.StatusOK, user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	filters := make(map[string]string)

	if r.FormValue("nickname") != "" {
		filters["nickname"] = r.FormValue("nickname")

	}

	if r.FormValue("email") != "" {
		filters["email"] = r.FormValue("email")

	}

	if r.FormValue("country") != "" {
		filters["country"] = r.FormValue("country")

	}

	users, err := RetrieveUsers(filters)

	if err != nil {
		log.Println(err.Error())

		utils.ReturnJsonResponse(
			w,
			http.StatusNotFound,
			map[string]string{"error": "The given filters don't seem valid."},
		)

		return
	}

	utils.ReturnJsonResponse(w, http.StatusOK, users)
}

func ValidateUser(u *User) (bool, error) {
	if u.IsValidNickname() != true {
		errorMessage := errors.New("The submitted nickname is not valid. It should have a length between 3 and 12 characters.")

		return false, errorMessage
	}

	if u.IsValidPassword() != true {
		errorMessage := errors.New("The submitted password is not valid. It should have a length of at least 8 characters.")

		return false, errorMessage
	}

	if u.IsValidEmail() != true {
		errorMessage := errors.New("The submitted Email is not valid.")

		return false, errorMessage
	}

	country := new(Country)
	country.IsoAlphaCode = u.Country

	_, countryError := country.IsValidCountry()

	if countryError != nil {
		countryError := errors.New("The submitted Country code for this user is invalid. Please use an ISO 3 characters code")
		log.Println(countryError.Error())

		return false, countryError
	}

	return true, nil
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
			map[string]string{"error": "Invalid resquest payload"},
		)

		return
	}

	_, validationError := ValidateUser(&u)

	if validationError != nil {
		log.Println(validationError.Error())

		utils.ReturnJsonResponse(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": validationError.Error(),
			},
		)
		return
	}

	err = InsertUser(&u)

	if err != nil {
		log.Println(err.Error())

		utils.ReturnJsonResponse(
			w,
			http.StatusInternalServerError,
			map[string]string{"status": "Unexpected error"},
		)

		return
	}

	statusMessage := "The user has successfully been created"

	err = events.PublishMessage(
		"CreatedUser",
		map[string]string{
			"status":       statusMessage,
			"userNickname": string(u.Nickname),
		},
	)

	if err != nil {
		log.Println(err.Error())
	}

	utils.ReturnJsonResponse(
		w,
		http.StatusOK,
		map[string]string{"status": statusMessage},
	)
}

func PutUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var userId string

	if Id, ok := vars["id"]; ok {
		userId = string(Id)
	} else {
		utils.ReturnJsonResponse(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "The provided user id is not valid."},
		)

		return
	}

	var user User

	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&user)

	defer r.Body.Close()

	if decodeErr != nil {
		log.Println(decodeErr.Error())
		utils.ReturnJsonResponse(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "Invalid resquest payload"},
		)

		return
	}

	_, validationError := ValidateUser(&user)

	if validationError != nil {
		log.Println(validationError.Error())

		utils.ReturnJsonResponse(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": validationError.Error(),
			},
		)
		return
	}

	err := UpdateUser(userId, &user)

	if err != nil {
		log.Println(err.Error())

		utils.ReturnJsonResponse(
			w,
			http.StatusInternalServerError,
			map[string]string{"status": "Unexpected error"},
		)

		return
	}

	statusMessage := "The user has successfully been updated"

	err = events.PublishMessage(
		"UpdatedUser",
		map[string]string{
			"status": statusMessage,
			"userId": userId,
		},
	)

	if err != nil {
		log.Println(err.Error())
	}

	utils.ReturnJsonResponse(
		w,
		http.StatusOK,
		map[string]string{"status": statusMessage},
	)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var userId string

	if Id, ok := vars["id"]; ok {
		userId = string(Id)
	} else {
		utils.ReturnJsonResponse(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "The provided user id is not valid or not given."},
		)

		return
	}

	err := RemoveUser(userId)

	if err != nil {
		log.Println(err.Error())

		utils.ReturnJsonResponse(
			w,
			http.StatusNotFound,
			map[string]string{"error": "The user does not exist."},
		)

		return
	}

	statusMessage := "The provided user has been successfully deleted."

	err = events.PublishMessage(
		"DeletedUser",
		map[string]string{
			"status": statusMessage,
			"userId": userId,
		},
	)

	if err != nil {
		log.Println(err.Error())
	}

	utils.ReturnJsonResponse(
		w,
		http.StatusOK,
		map[string]string{"status": statusMessage},
	)
}
