package users

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gnoirzox/faceit-users/utils"
)

func RetrieveUser(userId string) (*User, error) {
	db, err := utils.OpenDBConnection()
	defer db.Close()

	if err != nil {
		log.Println("%s: %s", "Could not connect to the database", err)

		return nil, err
	}

	row := db.QueryRow("SELECT firstname, lastname, nickname, email, country.name AS country FROM users INNER JOIN country ON users.country_code = country.alpha_code WHERE _id = $1", userId)

	var user User

	err = row.Scan(&user.Firstname, &user.Lastname, &user.Nickname, &user.Email, &user.Country)

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	return &user, nil
}

func RetrieveUsers(filters map[string]string) ([]User, error) {
	db, err := utils.OpenDBConnection()
	defer db.Close()

	if err != nil {
		log.Println("%s: %s", "Could not connect to the database", err)

		return nil, err
	}

	where := "WHERE 1=1"

	if nickname, ok := filters["nickname"]; ok {
		where += fmt.Sprintf(" AND nickname = '%s'", strings.Replace(nickname, "'", "", -1))
	}

	if email, ok := filters["email"]; ok {
		where += fmt.Sprintf(" AND email = '%s'", strings.Replace(email, "'", "", -1))
	}

	if country, ok := filters["country"]; ok {
		where += fmt.Sprintf(" AND country_code = '%s'", strings.Replace(country, "'", "", -1))
	}

	// TODO: need to improve it to avoid more clever SQL injections..
	queryString := "SELECT firstname, lastname, nickname, email, country.name AS country FROM users INNER JOIN country ON users.country_code = country.alpha_code " + where

	rows, err := db.Query(queryString)

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	users := make([]User, 0)

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.Firstname,
			&user.Lastname,
			&user.Nickname,
			&user.Email,
			&user.Country,
		)

		if err != nil {
			log.Println(err.Error())

			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func InsertUser(user *User) error {
	db, err := utils.OpenDBConnection()
	defer db.Close()

	if err != nil {
		log.Println("%s: %s", "Could not connect to the database", err)

		return err
	}

	insert, err := db.Prepare(
		"INSERT INTO users (firstname, lastname, nickname, password, email, country_code) VALUES ($1, $2, $3, $4, $5, $6)",
	)
	defer insert.Close()

	if err != nil {
		log.Println(err.Error())

		return err
	}

	hashedPassword, hashedPasswordError := utils.HashPassword(user.Password)

	if hashedPasswordError != nil {
		log.Println(hashedPasswordError.Error())

		return err
	}

	_, err = insert.Exec(
		&user.Firstname,
		&user.Lastname,
		&user.Nickname,
		hashedPassword,
		&user.Email,
		&user.Country,
	)

	if err != nil {
		log.Println(err.Error())

		return err
	}

	return nil
}

func UpdateUser(userId string, user *User) error {
	db, err := utils.OpenDBConnection()

	if err != nil {
		log.Println("%s: %s", "Could not connect to the database", err)

		return err
	}

	transaction, err := db.Begin()
	defer transaction.Rollback()

	if err != nil {
		log.Println(err.Error())

		return err
	}

	if userId != "" {
		_, err := RetrieveUser(userId)

		if err != nil {
			return err
		}
	} else {
		err := errors.New("No UserId provided, a user_id is mandatory to update a given user.")

		log.Println(err.Error())

		return err
	}

	update, err := transaction.Prepare(
		"UPDATE users SET firstname = $1, lastname = $2, nickname = $3, password = $4, email = $5, country_code = $6 WHERE _id = $7",
	)

	if err != nil {
		log.Println(err.Error())

		return err
	}

	defer update.Close()

	if err != nil {
		log.Println(err.Error())

		return err
	}

	hashedPassword, hashedPasswordError := utils.HashPassword(user.Password)

	if hashedPasswordError != nil {
		log.Println(hashedPasswordError.Error())

		return err
	}

	_, err = update.Exec(
		&user.Firstname,
		&user.Lastname,
		&user.Nickname,
		hashedPassword,
		&user.Email,
		&user.Country,
		userId,
	)

	if err != nil {
		log.Println(err.Error())

		return err
	}

	err = transaction.Commit()

	if err != nil {
		log.Println(err.Error())

		return err
	}

	return nil
}

func RemoveUser(UserId string) error {
	db, err := utils.OpenDBConnection()
	defer db.Close()

	if err != nil {
		log.Println("%s: %s", "Could not connect to the database", err)

		return err
	}

	delete, err := db.Prepare("DELETE FROM users WHERE _id = $1")
	defer delete.Close()

	_, err = delete.Exec(&UserId)

	if err != nil {
		log.Println(err.Error())

		return err
	}

	return nil
}
