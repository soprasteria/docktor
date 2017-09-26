package controllers

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/soprasteria/docktor/server/models"
	"github.com/soprasteria/docktor/server/modules/auth"
	"github.com/soprasteria/docktor/server/modules/users"
	"github.com/soprasteria/docktor/server/types"
	"gopkg.in/mgo.v2/bson"
)

// Users contains all group handlers
type Users struct {
}

//GetAll users from docktor
func (u *Users) GetAll(c echo.Context) error {
	docktorAPI := c.Get("api").(*models.Docktor)
	webservice := users.Rest{Docktor: docktorAPI}
	users, err := webservice.GetAllUserRest()
	if err != nil {
		log.WithError(err).Error("Unable to get all users")
		return c.String(http.StatusInternalServerError, "Unable to get all users because of technical error. Retry later.")
	}
	return c.JSON(http.StatusOK, users)
}

// Update user into docktor
// Only admin and current user is able to update a user
func (u *Users) Update(c echo.Context) error {
	docktorAPI := c.Get("api").(*models.Docktor)
	authenticatedUser, err := u.getUserFromToken(c)
	if err != nil {
		return c.String(http.StatusUnauthorized, auth.ErrInvalidCredentials.Error())
	}

	// Get User from body
	id := c.Param("userID")
	var userRest users.UserRest
	err = c.Bind(&userRest)
	if err != nil {
		log.WithError(err).Error("Unable to bind user to save")
		return c.String(http.StatusBadRequest, "Unable to parse user received from client")
	}

	// This route is only for update of existing user.
	// Another route exists for create a new user
	userRest.ID = id
	if userRest.ID == "" {
		log.Warn("Someone tried to use user updating route with empty ID")
		return c.String(http.StatusBadRequest, "Invalid user ID. User can not be created with this route. Please register.")
	}

	// Only admin or current user are authorized to modify user
	if authenticatedUser.ID != userRest.ID && !authenticatedUser.IsAdmin() {
		log.Errorf("User %v tried to update user %v but is not admin !", authenticatedUser.Username, userRest.ID)
		return c.String(http.StatusForbidden, ErrNotAuthorized.Error())
	}

	// Validate fields from validator tags for common types
	if err = c.Validate(userRest); err != nil {
		log.WithError(err).Warnf("User %s tried to save user %s, but some fields are not valid", authenticatedUser.ID, userRest.ID)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Some fields of user are not valid: %v", err))
	}

	var email, displayName, firstName, lastName *string
	var role *types.Role
	var tags, favorites []bson.ObjectId

	if authenticatedUser.IsAdmin() {
		log.WithFields(log.Fields{
			"username": userRest.Username,
			"newTags":  userRest.Tags,
			"newRole":  userRest.Role,
		}).Debug("Modifying user as Admin")
		// An admin is allowed to modify the following fields
		tags = userRest.Tags
		role = &userRest.Role
	}

	if authenticatedUser.ID == userRest.ID || authenticatedUser.IsAdmin() {
		// A user is allowed to modify the following fields from his own profile
		email = &userRest.Email
		displayName = &userRest.DisplayName
		firstName = &userRest.FirstName
		lastName = &userRest.LastName
		favorites = userRest.Favorites
	}

	// Automatically improve qualityf of data by keeping only existing external entities
	tags = existingTags(docktorAPI, tags)
	favorites = existingGroups(docktorAPI, favorites)

	webservice := users.Rest{Docktor: docktorAPI}
	res, err := webservice.UpdateUser(userRest.ID, email, displayName, firstName, lastName, role, tags, favorites)
	if err != nil {
		log.WithError(err).Errorf("Unable to update user %v because of unexpected error", userRest.ID)
		return c.String(http.StatusInternalServerError, "Unable to update user because of technical error. Retry later.")
	}

	return c.JSON(http.StatusOK, res)
}

//Delete user into docktor
func (u *Users) Delete(c echo.Context) error {
	docktorAPI := c.Get("api").(*models.Docktor)
	id := c.Param("userID")

	authenticatedUser, err := u.getUserFromToken(c)
	if err != nil {
		return c.String(http.StatusForbidden, auth.ErrInvalidCredentials.Error())
	}

	if authenticatedUser.ID != id && !authenticatedUser.IsAdmin() {
		// Admins can delete any users but user can only delete his own account
		log.Errorf("User %v tried to delete user %v but is not admin !", authenticatedUser.Username, id)
		return c.String(http.StatusForbidden, "You do not have rights to delete this user")
	}

	// Delete the user
	res, err := docktorAPI.Users().Delete(bson.ObjectIdHex(id))
	if err != nil {
		log.WithError(err).Errorf("Unable to delete user %v because of unexpected error", id)
		return c.String(http.StatusInternalServerError, "Unable to delete user because of technical error. Retry later.")
	}

	// Remove members on all groups as we delete it
	rmInfo, err := docktorAPI.Groups().RemoveMember(bson.ObjectIdHex(id))
	if err != nil {
		log.WithField("info", rmInfo).WithField("userId", id).Warn("Could not remove member from groups after deleting user")
	}

	return c.String(http.StatusOK, res.Hex())
}

// ChangePasswordOptions is a structure containing data used to change a password
// This struct will be unmarshalled from a HTTP request body
type ChangePasswordOptions struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// ChangePassword changes the password of a user
func (u *Users) ChangePassword(c echo.Context) error {

	var options ChangePasswordOptions
	err := c.Bind(&options)
	if err != nil {
		return c.String(http.StatusBadRequest, "Body not recognized")
	}

	authenticatedUser, err := u.getUserFromToken(c)
	if err != nil {
		return c.String(http.StatusForbidden, auth.ErrInvalidCredentials.Error())
	}

	id := c.Param("userID")
	if authenticatedUser.ID != id {
		log.Errorf("User %v tried to change password of user %v !", authenticatedUser.Username, id)
		return c.String(http.StatusForbidden, "Can't change password of someone else")
	}

	if options.NewPassword == "" || len(options.NewPassword) < 6 {
		log.Warn("User %v tried to change password that does not match security rules", authenticatedUser.ID)
		return c.String(http.StatusForbidden, "New password should not be empty and be at least 6 characters")
	}

	docktorAPI := c.Get("api").(*models.Docktor)
	webservice := auth.Authentication{Docktor: docktorAPI}
	err = webservice.ChangePassword(authenticatedUser.ID, options.OldPassword, options.NewPassword)

	if err != nil {
		if err == auth.ErrInvalidOldPassword {
			log.WithError(err).Warnf("User %v tried to change his password but invalid old password", authenticatedUser.Username)
			return c.String(http.StatusForbidden, err.Error())
		}
		log.WithError(err).Errorf("Unable to change password of user %v because of technical error", authenticatedUser.Username)
		return c.String(http.StatusInternalServerError, "Unable to change user password because of technical error. Retry later.")
	}

	return c.JSON(http.StatusOK, "")
}

// Profile returns the profile of the connecter user
func (u *Users) Profile(c echo.Context) error {
	user, err := u.getUserFromToken(c)
	if err != nil {
		return c.String(http.StatusUnauthorized, auth.ErrInvalidCredentials.Error())
	}

	return c.JSON(http.StatusOK, user)
}

//Get user from docktor
func (u *Users) Get(c echo.Context) error {
	// No access control on purpose
	user := c.Get("user").(users.UserRest)
	return c.JSON(http.StatusOK, user)
}

func (u *Users) getUserFromToken(c echo.Context) (users.UserRest, error) {
	docktorAPI := c.Get("api").(*models.Docktor)
	userToken := c.Get("user-token").(*jwt.Token)

	claims := userToken.Claims.(*auth.MyCustomClaims)

	webservice := users.Rest{Docktor: docktorAPI}
	return webservice.GetUserRest(claims.Username)
}
