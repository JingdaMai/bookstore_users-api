package users

import (
	"fmt"
	"github.com/JingdaMai/bookstore_items-api/datasources/postgresql/users_db"
	"github.com/JingdaMai/bookstore_items-api/utils/date_utils"
	"github.com/JingdaMai/bookstore_items-api/utils/errors"
	"strings"
)

const (
	queryInsertUser = "INSERT INTO users(first_name,last_name,email,date_created) VALUES ($1,$2,$3,$4) RETURNING id;"
)

var (
	userDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {
	if err := users_db.Client.Ping(); err != nil {
		panic(err)
	}

	result := userDB[user.Id]
	if result == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
	}
	user.Id = result.Id
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.DateCreated = result.DateCreated

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
	err = stmt.QueryRow(user.FirstName, user.LastName, user.Email, user.DateCreated).Scan(&userId)
	if err != nil {
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_email_key"`) {
			return errors.NewBadRequestError(fmt.Sprintf("email %s already exists", user.Email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save users: %s", err.Error()))
	}

	user.Id = userId

	return nil
}
