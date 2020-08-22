package users

import (
	"github.com/JingdaMai/bookstore_items-api/domain/users"
	"github.com/JingdaMai/bookstore_items-api/services"
	"github.com/JingdaMai/bookstore_items-api/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var user users.User

	// parse request body
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	// create user in the service
	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}

	// return result
	c.JSON(http.StatusCreated, result)
}

func GetUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "not implemented")
}

func SearchUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "not implemented")
}
