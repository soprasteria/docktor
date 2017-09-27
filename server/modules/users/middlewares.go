package users

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/soprasteria/docktor/server/storage"
)

// RetrieveUser find user using id param and put it in echo.Context
func RetrieveUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		docktorAPI := c.Get("api").(*storage.Docktor)
		id := c.Param("userID")
		if id == "" {
			return c.String(http.StatusBadRequest, UserInvalidID)
		}
		user, err := docktorAPI.Users().FindByID(id)
		if err != nil {
			return c.String(http.StatusNotFound, fmt.Sprintf(UserNotFound, id))
		}

		c.Set("user", GetUserRest(user))
		return next(c)
	}
}
