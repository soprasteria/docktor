package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/soprasteria/docktor/server/modules/auth"
	"github.com/soprasteria/docktor/server/modules/users"
	"github.com/soprasteria/docktor/server/storage"
)

// ErrNotAuthorized is an error when someone is trying to access unauthorized ressource for a given role
var ErrNotAuthorized = errors.New("You're trying to access unauthorized ressources for your role")

// Auth contains all login handlers
type Auth struct {
}

// Token is a JWT Token
type Token struct {
	ID   string         `json:"id_token,omitempty"`
	User users.UserRest `json:"user,omitempty"`
}

func newAuthAPI(c echo.Context) auth.Authentication {
	// Handle APIs from Echo context
	docktorAPI := c.Get("api").(*storage.Docktor)
	ldapAPI := c.Get("ldap")
	var ldap *auth.LDAP
	if ldapAPI != nil {
		ldap = ldapAPI.(*auth.LDAP)
	}
	return auth.Authentication{
		Docktor: docktorAPI,
		LDAP:    ldap,
	}
}

// Register create an account
func (a *Auth) Register(c echo.Context) error {

	// Get form data
	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")
	firstname := c.FormValue("firstname")
	lastname := c.FormValue("lastname")

	// Check form data
	registerQuery := auth.RegisterUserQuery{
		Username:  username,
		Password:  password,
		Email:     email,
		Firstname: firstname,
		Lastname:  lastname,
	}
	if err := c.Validate(registerQuery); err != nil {
		log.WithError(err).Warnf("Parameters to register are not valid for user %v", username)
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Handle APIs from Echo context
	login := newAuthAPI(c)

	// Log in the application
	err := login.RegisterUser(&registerQuery)
	if err != nil {
		if err == auth.ErrUsernameAlreadyTaken {
			log.WithError(err).Warnf("Someone tried to register with an existing username %v", username)
			return c.String(http.StatusForbidden, auth.ErrUsernameAlreadyTaken.Error())
		}
		log.WithError(err).WithField("username", username).Error("Unable to register user")
		return c.String(http.StatusInternalServerError, "Unable to register user because of technical error. Retry later.")
	}

	// Generates a valid token
	token, err := login.CreateLoginToken(username)
	if err != nil {
		log.WithError(err).WithField("username", username).Error("Token creation failed while registering user")
		return c.String(http.StatusInternalServerError, "Unable to register user because of technical error. Retry later.")
	}

	// Get the user from database
	webservice := users.Rest{Docktor: login.Docktor}
	user, err := webservice.GetUserRest(username)
	if err != nil {
		log.WithError(err).WithField("username", username).Error("User retrieval failed while registering user")
		return c.String(http.StatusInternalServerError, "Unable to register user because of technical error. Retry later.")
	}

	return c.JSON(http.StatusOK, Token{ID: token, User: user})

}

//Login handles the login of a user
//When user is authorized, it creates a JWT Token https://jwt.io/introduction/ that will be store on client
func (a *Auth) Login(c echo.Context) error {
	// Get input parameters
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Check input parameters
	if username == "" {
		log.Warn("Someone tried to login with empty username")
		return c.String(http.StatusForbidden, "Username should not be empty")
	}

	if len(password) < 6 {
		log.Warn("Someone tried to login with a password that does not match security rules")
		return c.String(http.StatusForbidden, "Password length should be at least 6 characters")
	}

	// Handle APIs from Echo context
	login := newAuthAPI(c)

	// Log in the application
	err := login.AuthenticateUser(&auth.LoginUserQuery{
		Username: username,
		Password: password,
	})
	if err != nil {
		if err == auth.ErrInvalidCredentials {
			log.WithError(err).Warnf("Someone tried to login with wrong credentials: %v", username)
			return c.String(http.StatusForbidden, auth.ErrInvalidCredentials.Error())
		}
		log.WithError(err).Error("User authentication failed")
		return c.String(http.StatusInternalServerError, "Unable to login user because of technical error. Retry later.")
	}

	// Generates a valid token
	token, err := login.CreateLoginToken(username)
	if err != nil {
		log.WithError(err).WithField("username", username).Error("Login token creation failed while logging in user")
		return c.String(http.StatusInternalServerError, "Unable to login user because of technical error. Retry later.")
	}

	// Get the user from database
	webservice := users.Rest{Docktor: login.Docktor}
	user, err := webservice.GetUserRest(username)
	if err != nil {
		log.WithError(err).WithField("username", username).Error("Unable to retrieve user while logging in user")
		return c.String(http.StatusInternalServerError, "Unable to login user because of technical error. Retry later.")
	}

	return c.JSON(http.StatusOK, Token{ID: token, User: user})
}

//ResetPassword handles resets the password of someone
//When user reset his password, it generates an email to this person with a link to change his.
// The password is not reset in database because someone with bad intentions could use this feature to prevent someone else to login.
func (a *Auth) ResetPassword(c echo.Context) error {
	// Get input parameters

	username := c.FormValue("username")

	// Check input parameters
	if username == "" {
		log.Warn("Someone tried to reset password with empty username")
		return c.String(http.StatusForbidden, "Username should not be empty")
	}

	// Handle APIs from Echo context
	login := newAuthAPI(c)

	// Check whether the password can be reset or not
	user, err := login.IsPasswordUserResetable(username)
	if err != nil {
		log.WithError(err).Warnf("Someone tried to reset password of user %v because was not authorized to do so.", username)
		return c.String(http.StatusForbidden, err.Error())
	}
	// Create a token that will be used for URL
	token, err := login.CreateResetPasswordToken(user.Username)
	if err != nil {
		log.WithError(err).Errorf("Unable to generate the password change URL for user %v", username)
		return c.String(http.StatusInternalServerError, "Unable to reset password because of technical error. Retry Later.")
	}

	var protocol = "http"
	if c.Request().TLS != nil {
		protocol = "https"
	}

	url := fmt.Sprintf("%s://%s/change_reset_password?token=%s", protocol, c.Request().Host, token)
	go auth.SendResetPasswordEmail(user, url)

	return c.JSON(http.StatusOK, "OK")
}

//ChangeResetPassword changes the password of someone that reset it before
//When user changes the password, he's automatically connected
func (a *Auth) ChangeResetPassword(c echo.Context) error {
	// Get input parameters

	newPassword := c.FormValue("newPassword")
	token := c.FormValue("token")

	// Check input parameters
	if newPassword == "" {
		log.Warn("Someone tried to reset change a resetted password with empty new password")
		return c.String(http.StatusForbidden, "NewPassword should not be empty")
	}

	if token == "" {
		log.Warn("Someone tried to reset change a resetted password with empty reset token")
		return c.String(http.StatusForbidden, "Token should not be empty")
	}

	// Handle APIs from Echo context
	login := newAuthAPI(c)

	// Change the password of the user in DB
	user, err := login.ChangeResetPasswordUser(token, newPassword)
	if err != nil {
		log.WithError(err).Error("Unable to reset password in database")
		return c.String(http.StatusForbidden, err.Error())
	}

	// Generates a valid token
	authenticationToken, err := login.CreateLoginToken(user.Username)
	if err != nil {
		log.WithError(err).WithField("username", user.Username).Error("Login token creation failed after a password change")
		return c.String(http.StatusInternalServerError, "Unable to change reset password because of technical error. Retry Later.")
	}

	// Get the user from database
	webservice := users.Rest{Docktor: login.Docktor}
	userRest, err := webservice.GetUserRest(user.Username)
	if err != nil {
		log.WithError(err).WithField("username", user.Username).Error("Unable to fetch user after password change")
		return c.String(http.StatusInternalServerError, "Unable to change reset password because of technical error. Retry Later.")
	}

	return c.JSON(http.StatusOK, Token{ID: authenticationToken, User: userRest})
}
