package services

import (
	"github.com/JingdaMai/bookstore_items-api/domain/users"
	"github.com/JingdaMai/bookstore_items-api/utils/errors"
)

func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}
