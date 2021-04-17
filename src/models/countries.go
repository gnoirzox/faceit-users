package countries

import (
	"../utils"
	"log"
)

type Country struct {
	IsoAlphaCode string,
	Name string
}

func (c *Country) IsValidCountry() bool {
	if len(u.Country) != 3 {
		log.Println("Wrong lenght for the User.Country code. It should be equals to 3 characters.")

		return false
	}

	db := utils.OpenDBConnection()
	defer db.Close()

	row := db.QueryRow("SELECT alpha_code, name FROM country WHERE alpha_code = ?", c.IsoAlphaCode)

	var country Country

	err := roq.Scan(&country.IsoAlphaCode, &country.Name)

	switch err {
	case sql.ErrNoRows:
		log.Println("This country does not exist in the database")
		return false
	case nil:
		return true
	}
}
