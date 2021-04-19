package users

import (
	"../utils"
	"database/sql"
	"log"
)

type Country struct {
	IsoAlphaCode string
	Name         string
}

func (c *Country) IsValidCountry() bool {
	if len(c.IsoAlphaCode) != 3 {
		log.Println("Wrong lenght for the User.Country code. It should be equals to 3 characters.")

		return false
	}

	db, err := utils.OpenDBConnection()

	if err != nil {
		log.Println("%s: %s", "Could not connect to the database", err)

		return false
	}

	defer db.Close()

	row := db.QueryRow("SELECT alpha_code, name FROM country WHERE alpha_code = ?", c.IsoAlphaCode)

	var country Country

	err = row.Scan(&country.IsoAlphaCode, &country.Name)

	switch err {
	case sql.ErrNoRows:
		log.Println("This country does not exist in the database")
		return false
	}

	return true
}
