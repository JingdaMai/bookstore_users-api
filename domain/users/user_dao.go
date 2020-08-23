package users

import (
	"github.com/JingdaMai/bookstore_items-api/datasources/postgresql/users_db"
	"github.com/JingdaMai/bookstore_items-api/utils/date_utils"
	"github.com/JingdaMai/bookstore_items-api/utils/errors"
	"github.com/JingdaMai/bookstore_items-api/utils/postgres_utils"
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES ($1, $2, $3, $4) RETURNING id;"
	queryGetUser    = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id=$1;"
	queryUpdateUser = "UPDATE users SET first_name=$1, last_name=$2, email=$3 WHERE id=$4;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); getErr != nil {
		return postgres_utils.ParseError(getErr)
	}

	return nil
}

func (user *User) Save() *errors.RestErr {
	// prepare INSERT statement
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	user.DateCreated = date_utils.GetNowString()

	var userId int64

	// execute INSERT statement
	saveErr := stmt.QueryRow(user.FirstName, user.LastName, user.Email, user.DateCreated).Scan(&userId)
	if saveErr != nil {
		return postgres_utils.ParseError(saveErr)
	}

	user.Id = userId

	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		return postgres_utils.ParseError(err)
	}

	return nil
}
