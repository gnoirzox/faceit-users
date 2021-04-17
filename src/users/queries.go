package users

import (
	"container/list"
	"errors"
	"log"

	"../utils"
)

func RetrieveUser(UserId string) (*User, error) {
	db := utils.OpenDBConnection()
	defer db.Close()

	row := db.QueryRow("SELECT firstname, lastname, nickname, email, country_code FROM users WHERE _id = ?", UserId)

	var user User

	err := row.Scan(&user.Firstname, &user.Lastname, &user.Nickname, &user.Email, &user.Country, &UserId)

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	return &user, nil
}

func RetrieveUsers(filters map[string]string) (*list.List, error) {
	db := utils.OpenDBConnection()
	defer db.Close()

	where := "WHERE 1=1"

	if nickname, ok := filters["nickname"]; ok {
		where += " AND nickname = " + nickname
	}

	if email, ok := filters["email"]; ok {
		where += " AND email = " + email
	}

	if country, ok := filters["country"]; ok {
		where += " AND country_code = " + country
	}

	queryString := "SELECT firstname, lastname, nickname, email, country_code FROM users " + where

	rows, err := db.Query(queryString)

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	users := list.New()

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

		users.PushBack(user)
	}

	return users, nil
}

func InsertUser(user *User) error {
	db := utils.OpenDBConnection()
	defer db.Close()

	insert, err := db.Prepare(
		"INSERT INTO users (firstname, lastname, nickname, password, email, country_code) VALUES (?, ?, ?, ?, ?, ?)",
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
	db := utils.OpenDBConnection()

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
		"UPDATE users SET firstname = ?, lastname = ?, nickname = ?, password = ?, email = ?, country_code = ? WHERE _id = ?",
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
	db := utils.OpenDBConnection()
	defer db.Close()

	delete, err := db.Prepare("DELETE users WHERE _id = ?")
	defer delete.Close()

	_, err = delete.Exec(&UserId)

	if err != nil {
		log.Println(err.Error())

		return err
	}

	return nil
}
