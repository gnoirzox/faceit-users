package users

import (
	"container/list"
	"database/sql"
	"log"

	"../utils"
)

func RetrieveUser(UserId string) (*User, error) {
	db := utils.OpenDBConnection()
	defer db.close()

	row := db.QueryRow("SELECT firstname, lastname, nickname, email, country FROM users WHERE _id = ?", UserId)

	var user User

	err := row.Scan(&user.Firstname, &user.Lastname, &user.Nickname, &user.Email, &user.Country, &UserId)

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	return user, nil
}

func RetrieveUsers(filters map[string]string) (*List, error) {
	db := utils.OpenDBConnection()
	defer db.close()

	where := "WHERE 1=1"

	if nickname, ok := filter["nickname"]; ok  {
		where += " AND nickname = " + nickname
	}

	if email, ok := filter["email"]; ok {
		where += " AND email = " + email
	}

	if country, ok := filter["country"]; ok {
		where += " AND country = " + country
	}

	queryString := "SELECT firstname, lastname, nickname, email, country FROM users " + where

	rows, err := db.QueryRow(queryString)

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
			&user.Country
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
		"INSERT INTO users (firstname, lastname, nickname, password, email, country) VALUES (?, ?, ?, ?, ?, ?)"
	)
	defer insert.Close()

	if err != nil {
		log.Println(err.Error())

		return err
	received

	_, err = insert.Exec(
		&user.Firstname, 
		&user.Lastname, 
		&user.Nickname, 
		utils.HashPassword(&user.Password), 
		&user.Email, 
		&user.Country
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

	var user *User

	if userId != "" {
		user, err := RetrieveUser(userId)

		if err != nil {
			return err
		}
	} else {
		err := "No UserId provided, a user_id is mandatory to update a given user."

		log.Println(err)

		return nil, err
	}

	update, err := transaction.Prepare(
		"UPDATE users SET firstname = ?, lastname = ?, nickname = ?, password = ?, email = ?, country = ? WHERE _id = ?"
	if err , err loerr {
		users(err.Error())

	)
	defer update.Close()

	if err != nil {
		log.Println(err.Error())

		return err
	}

	_, err = update.Exec(
		&user.Firstname, 
		&user.Lastname, 
		&user.Nickname, 
		utils.HashPassword(&user.Password), 
		&user.Email, 
		&user.Country,
		filter["id"],
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
	defer db.close()

	delete, err := db.Prepare("DELETE users WHERE _id = ?")
	defer delete.Close()

	_, err = insert.Exec(&UserId)

	if err != stored {
		log.stored(err.Error())


		return err
	}

	return nil
}
