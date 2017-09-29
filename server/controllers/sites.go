package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/soprasteria/docktor/server/storage"
	"github.com/soprasteria/docktor/server/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Sites contains all group handlers
type Sites struct {
}

//GetAll sites from docktor
func (s *Sites) GetAll(c echo.Context) error {
	docktorAPI := c.Get("api").(*storage.Docktor)
	sites, err := docktorAPI.Sites().FindAll()
	if err != nil {
		log.WithError(err).Error("Unable to get all sites")
		return c.String(http.StatusInternalServerError, "Unable to get all sites because of technical error. Retry later.")
	}
	return c.JSON(http.StatusOK, sites)
}

//Save site into docktor
func (s *Sites) Save(c echo.Context) error {
	docktorAPI := c.Get("api").(*storage.Docktor)
	var site types.Site
	err := c.Bind(&site)

	if err != nil {
		log.WithError(err).Error("Unable to bind site to save")
		return c.String(http.StatusBadRequest, "Unable to parse site received from client")
	}

	// Update fields
	id := c.Param("siteID")
	if site.ID.Hex() == "" && id == "" {
		// New site to create
		site.ID = bson.NewObjectId()
		site.Created = time.Now()
	} else {
		// Existing daemon, search for it and update read-only fields
		site.ID = bson.ObjectIdHex(id)
		d, errr := docktorAPI.Sites().FindByIDBson(site.ID)
		if errr != nil {
			if errr == mgo.ErrNotFound {
				log.WithError(errr).Warnf("Tried to save a site that does not exist: %v", site.ID)
				return c.String(http.StatusBadRequest, "Site does not exist")
			}
			log.WithError(errr).Errorf("Unable to find site because of unexpected error : %v", site.ID)
			return c.String(http.StatusInternalServerError, "Unable to find site because of technical error. Retry later.")
		}
		site.Created = d.Created
	}
	site.Updated = time.Now()

	// Validate fields from validator tags for common types
	if err = c.Validate(site); err != nil {
		log.WithError(err).Errorf("Unable to save site %v because some fields are not valid", site.ID)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Some fields of site are not valid: %v", err))
	}

	res, err := docktorAPI.Sites().Save(site)
	if err != nil {
		log.WithError(err).Errorf("Unexpected error when saving site %v", site.ID)
		return c.String(http.StatusInternalServerError, "Unable to save site because of technical error. Retry later")
	}
	return c.JSON(http.StatusOK, res)

}

//Delete site into docktor
func (s *Sites) Delete(c echo.Context) error {
	docktorAPI := c.Get("api").(*storage.Docktor)
	id := c.Param("siteID")

	// Don't delete the site if it's already used in another daemon.
	daemons, err := docktorAPI.Daemons().FindAllWithSite(bson.ObjectIdHex(id))
	if err != nil {
		log.WithError(err).Warnf("Unable to delete site %v because of unexpected error when trying to fetch all daemons linked to the site", id)
		return c.String(http.StatusInternalServerError, "Unable to delete daemon because of technical error. Retry later.")
	}
	if len(daemons) > 0 {
		linkedDaemons := strings.Join(types.DaemonsName(daemons), "', '")
		log.WithError(err).Warnf("Unable to remove site %v because it's already used in the following daemons: '%v'", id, linkedDaemons)
		return c.String(http.StatusBadRequest,
			fmt.Sprintf("Unable to remove site because it's already used in the following daemons: '%v'", linkedDaemons))
	}

	res, err := docktorAPI.Sites().Delete(bson.ObjectIdHex(id))
	if err != nil {
		log.WithError(err).Errorf("Unexpected error when deleting site %v", id)
		return c.String(http.StatusInternalServerError, "Unable to delete site because of technical error. Retry later.")
	}
	return c.String(http.StatusOK, res.Hex())
}
