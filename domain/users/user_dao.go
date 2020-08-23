package users

import (
	"fmt"
	"github.com/JingdaMai/bookstore_items-api/datasources/postgresql/users_db"
	"github.com/JingdaMai/bookstore_items-api/logger"
	"github.com/JingdaMai/bookstore_items-api/utils/errors"
	"github.com/JingdaMai/bookstore_items-api/utils/postgres_utils"
)

const (
	queryInsertUser       = "INSERT INTO users(first_name, last_name, email, date_created, password, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"
	queryGetUser          = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=$1;"
	queryUpdateUser       = "UPDATE users SET first_name=$1, last_name=$2, email=$3 WHERE id=$4;"
	queryDeleteUser       = "DELETE FROM users WHERE id=$1;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created FROM users WHERE status=$1;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		logger.Error("error when trying to get user by id", err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) Save() *errors.RestErr {
	// prepare INSERT statement
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	var userId int64

	// execute INSERT statement
	saveErr := stmt.QueryRow(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Password, user.Status).Scan(&userId)
	if saveErr != nil {
		logger.Error("error when trying to save user", err)
		return postgres_utils.ParseError(saveErr)
	}

	user.Id = userId

	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return postgres_utils.ParseError(err)
	}

	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when trying to delete user", err)
		return postgres_utils.ParseError(err)
	}

	return nil
}

func (user *User) FindByStatus(status string) (Users, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find user by status statement", err)
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find users by status", err)
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	var results Users

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); err != nil {
			logger.Error("error when scanning user row to user struct", err)
			return nil, postgres_utils.ParseError(err)
		}
		user.Status = status
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}

	return results, nil
}
