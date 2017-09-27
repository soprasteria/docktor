package controllers

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/soprasteria/docktor/server/modules/auth"
	"github.com/soprasteria/docktor/server/modules/daemons"
	"github.com/soprasteria/docktor/server/storage"
	"github.com/soprasteria/docktor/server/types"
	"github.com/soprasteria/docktor/server/utils"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Daemons contains all daemons handlers
type Daemons struct {
}

// GetAll daemons from docktor
func (d *Daemons) GetAll(c echo.Context) error {
	docktorAPI := c.Get("api").(*storage.Docktor)
	daemons, err := docktorAPI.Daemons().FindAll()
	if err != nil {
		log.WithError(err).Error("Unable to get all daemons")
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to get all daemons because of technical error: %v. Retry later.", err))
	}
	return c.JSON(http.StatusOK, daemons)
}

//Save daemon into docktor
func (d *Daemons) Save(c echo.Context) error {
	docktorAPI := c.Get("api").(*storage.Docktor)
	var daemon types.Daemon
	if err := c.Bind(&daemon); err != nil {
		log.WithError(err).Error("Unable to bind daemon to save")
		return c.String(http.StatusBadRequest, "Unable to parse daemon received from client")
	}

	// Update fields
	id := c.Param("daemonID")
	if id == "" {
		// New daemon to create
		daemon.ID = bson.NewObjectId()
		daemon.Created = time.Now()
	} else {
		// Existing daemon, search for it and update read-only fields
		daemon.ID = bson.ObjectIdHex(id)
		d, err := docktorAPI.Daemons().FindByIDBson(daemon.ID)
		if err != nil {
			if err == mgo.ErrNotFound {
				log.WithError(err).Warnf("Tried to save a daemon that does not exist: %v", daemon.ID)
				return c.String(http.StatusBadRequest, "Daemon does not exist")
			}
			log.WithError(err).Errorf("Unable to find daemon because of unexpected error : %v", daemon.ID)
			return c.String(http.StatusInternalServerError, "Unable to find daemon because of technical error. Retry later.")
		}
		daemon.Created = d.Created
	}
	if daemon.Protocol == types.HTTPProtocol {
		daemon.Ca = ""
		daemon.Key = ""
		daemon.Cert = ""
	}
	daemon.Updated = time.Now()

	// Validate fields from validator tags for common types
	if err := c.Validate(daemon); err != nil {
		log.WithError(err).Errorf("Unable to save daemon %v because some fields are not valid", daemon.ID)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Some fields of daemon are not valid: %v", err))
	}

	// Validate fields that cannot be validated by validator engine
	if err := daemon.Validate(); err != nil {
		log.WithError(err).Errorf("Unable to save daemon %v because some fields are not valid", daemon.ID)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Some fields of daemon are not valid: %v", err))
	}

	// Check that daemon site exists
	if _, err := docktorAPI.Sites().FindByIDBson(daemon.Site); err != nil {
		loginfo := log.Fields{
			"site":   daemon.Site,
			"daemon": daemon.ID,
		}
		if err == mgo.ErrNotFound {
			log.WithError(err).WithFields(loginfo).Warnf("Tried to save a daemon with given site but site does not exist")
			return c.String(http.StatusBadRequest, "Site does not exist")
		}
		log.WithError(err).WithFields(loginfo).Errorf("Tried to save a daemon with given site but unexpected error when fetching site")
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to check if site exist: %v. Retry later.", err))
	}

	// Keep only existing tags
	daemon.Tags = existingTags(docktorAPI, daemon.Tags)

	res, err := docktorAPI.Daemons().Save(daemon)
	if err != nil {
		log.WithError(err).Errorf("Unexpected error when saving daemon %v", daemon.ID)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to save daemon because of technical error: %v. Retry later.", err))
	}
	return c.JSON(http.StatusOK, res)
}

//Delete daemon into docktor
func (d *Daemons) Delete(c echo.Context) error {
	docktorAPI := c.Get("api").(*storage.Docktor)
	id := c.Param("daemonID")

	// TODO: return error when daemon is already used in another service/container
	// If it's used in a filesystem in a group, it's OK, it should be deleted when saving the group

	res, err := docktorAPI.Daemons().Delete(bson.ObjectIdHex(id))
	if err != nil {
		log.WithError(err).Errorf("Unexpected error when deleting daemon %v", id)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to delete daemon because of technical error: %v. Retry later.", err))
	}
	return c.String(http.StatusOK, res.Hex())
}

//Get daemon from docktor
func (d *Daemons) Get(c echo.Context) error {
	daemon := c.Get("daemon").(types.Daemon)
	authenticatedUser, err := getUserFromToken(c)
	if err != nil {
		return c.String(http.StatusForbidden, auth.ErrInvalidCredentials.Error())
	}
	if !authenticatedUser.IsAdmin() {
		// Fetch daemon, amputed of its sensible data when user is not admin
		return c.JSON(http.StatusOK, daemons.GetDaemonRest(daemon))
	}
	return c.JSON(http.StatusOK, daemon)
}

// GetInfo : get infos about daemon from docker
func (d *Daemons) GetInfo(c echo.Context) error {
	daemon := c.Get("daemon").(types.Daemon)
	forceParam := c.QueryParam("force")
	redisClient := utils.GetRedis(c)

	infos, err := daemons.GetInfo(daemon, redisClient, forceParam == "true")
	if err != nil {
		log.WithError(err).Errorf("Unexpected error when getting info/status from daemon %v", daemon.ID)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to get daemon info because of technical error: %v. Retry later.", err))
	}
	return c.JSON(http.StatusOK, infos)
}
