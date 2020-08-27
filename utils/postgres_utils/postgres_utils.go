package postgres_utils

import (
	"github.com/JingdaMai/bookstore_items-api/utils/errors"
	"github.com/lib/pq"
	"strings"
)

const ErrorNoRows = "no rows in result set"

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*pq.Error)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRows) {
			return errors.NewNotFoundError("no record matching given id")
		}
		return errors.NewInternalServerError("error parsing database response")
	}
	switch sqlErr.Code {
	case "23505":
		return errors.NewBadRequestError("duplicated key")
	}
	return errors.NewInternalServerError("error processing request")
}
