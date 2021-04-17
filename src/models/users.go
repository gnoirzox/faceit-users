package users

import (
	"log"
	"net/mail"
)

type User struct {
	Firstname string `json:"firstname,omitempty"`,
	Lastname string `json:"lastname,omitempty"`,
	Nickname string `json:"nickname,omitempty"`,
	Password string `json:"password"`,
	Email string `json:"email,omitempty"`,
	Country string `json:"country,omitempty"`,
}

func (u *User) IsValidNickname() bool {
	if len(u.Nickname) < 3 || len(u.Nickname) > 12 {
		log.Println("Wrong lenght for the User.Nickname. It should be between 3 and 12 characters.")

		return false
	}

	return true
}

func (u *User) IsValidPassword() bool {
	if len(u.Password) >= 8 {
		log.Println("The User.Password is too short. It should be at least 8 characters.")

		return false
	}

	return true
}

func (u *User) IsValidEmail() bool {
	_, err := mail.ParseAddress(u.Email)

	if err != nil {
		log.Println("The provided email address is not respecting the RFC 5322 format for User.Email")
		log.Println(err.Error())

		return false
	}

	return true
}
