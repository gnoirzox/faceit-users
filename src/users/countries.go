package users

import (
	"database/sql"
	"errors"
	"github.com/gnoirzox/faceit-users/utils"
	"log"
	"strings"
)

type Country struct {
	IsoAlphaCode string
	Name         string
}

func (c *Country) IsValidCountry() (bool, error) {
	if len(c.IsoAlphaCode) != 3 {
		errorMessage := "Wrong lenght for the User.Country code. It should be equals to 3 characters."
		log.Println(errorMessage)

		err := errors.New(errorMessage)

		return false, err
	}

	db, err := utils.OpenDBConnection()

	if err != nil {
		log.Printf("%s: %s", "Could not connect to the database", err)

		return false, err
	}

	defer db.Close()

	row := db.QueryRow(
		"SELECT alpha_code, name FROM country WHERE alpha_code = $1",
		strings.ToUpper(c.IsoAlphaCode),
	)

	var country Country

	err = row.Scan(&country.IsoAlphaCode, &country.Name)

	switch {
	case err == sql.ErrNoRows:
		errorMessage := "This country does not exist in the database"
		log.Println(errorMessage)

		err = errors.New(errorMessage)

		return false, err
	case err != nil:
		log.Println(err.Error())
		return false, err
	}

	return true, nil
}
