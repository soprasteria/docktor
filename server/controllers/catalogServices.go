package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/soprasteria/docktor/server/storage"
	"github.com/soprasteria/docktor/server/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CatalogServices contains all services handlers
type CatalogServices struct {
}

// GetAll catalogServices from docktor
func (s *CatalogServices) GetAll(c echo.Context) error {
	docktorAPI := c.Get("api").(*storage.Docktor)
	catalogServices, err := docktorAPI.CatalogServices().FindAll()
	if err != nil {
		log.WithError(err).Error("Unable to get all catalogServices")
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to get all daemons because of technical error: %v. Retry later.", err))
	}
	return c.JSON(http.StatusOK, catalogServices)
}

// Save catalogService into docktor
func (s *CatalogServices) Save(c echo.Context) error {
	docktorAPI := c.Get("api").(*storage.Docktor)
	var catalogService types.CatalogService
	if err := c.Bind(&catalogService); err != nil {
		log.WithError(err).Error("Unable to bind catalogService to save")
		return c.String(http.StatusBadRequest, "Unable to parse catalogService received from client")
	}

	// If the ID is empty, it's a creation, so generate an object ID
	id := c.Param("catalogServiceID")
	if id == "" {
		// New catalogService to create
		catalogService.ID = bson.NewObjectId()
		catalogService.Created = time.Now()
	} else {
		// Existing catalogService
		catalogService.ID = bson.ObjectIdHex(id)
		s, err := docktorAPI.CatalogServices().FindByIDBson(catalogService.ID)
		if err != nil {
			if err == mgo.ErrNotFound {
				log.WithError(err).Warnf("Tried to save a catalogService that does not exist: %v", catalogService.ID)
				return c.String(http.StatusBadRequest, "CatalogService does not exist")
			}
			log.WithError(err).Errorf("Unable to find catalogService because of unexpected error : %v", catalogService.ID)
			return c.String(http.StatusInternalServerError, "Unable to find catalogService because of technical error. Retry later.")
		}
		catalogService.Created = s.Created
	}

	// Validate fields from validator tags for common types
	if err := c.Validate(catalogService); err != nil {
		log.WithError(err).Errorf("Unable to save catalogService %v because some fields are not valid", catalogService.ID)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Some fields of catalogService are not valid: %v", err))
	}

	// Validate fields that cIannot be validated by validator engine
	if err := catalogService.Validate(); err != nil {
		log.WithError(err).Errorf("Unable to save catalogService %v because some fields are not valid", catalogService.ID)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Some fields of catalogService are not valid: %v", err))
	}

	// Keep only existing and remove duplicates of external collections
	// Used to clean the service of old data when saving
	catalogService.Tags = existingTags(docktorAPI, catalogService.Tags)
	catalogService.Updated = time.Now()

	res, err := docktorAPI.CatalogServices().Save(catalogService)
	if err != nil {
		log.WithError(err).Errorf("Unexpected error when saving catalogService %v", catalogService.ID)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to save catalogService because of technical error: %v. Retry later.", err))
	}
	return c.JSON(http.StatusOK, res)
}

// existingCatalogServices return catalogServices filtered by existing ones
// Checks wether the Services actually exists in database
func existingCatalogServices(docktorAPI *storage.Docktor, catalogServicesIDs []bson.ObjectId) []bson.ObjectId {

	existingCatalogServiceIDs := []bson.ObjectId{}

	// Get all real groups
	existingcatalogServices, _ := docktorAPI.CatalogServices().FindAllByIDs(catalogServicesIDs)

	// Get their ids only
	for _, s := range existingcatalogServices {
		existingCatalogServiceIDs = append(existingCatalogServiceIDs, s.ID)
	}

	return existingCatalogServiceIDs
}

// Delete catalogServices into docktor
func (s *CatalogServices) Delete(c echo.Context) error {
	docktorAPI := c.Get("api").(*storage.Docktor)
	id := c.Param("catalogServiceID")

	res, err := docktorAPI.CatalogServices().Delete(bson.ObjectIdHex(id))
	if err != nil {
		log.WithError(err).Errorf("Unexpected error when deleting catalogService %v", id)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to delete catalogService because of technical error: %v. Retry later.", err))
	}
	return c.String(http.StatusOK, res.Hex())
}

// Get get a catalogService from docktor by its ID
func (s *CatalogServices) Get(c echo.Context) error {
	catalogService := c.Get("catalogService").(types.CatalogService)
	return c.JSON(http.StatusOK, catalogService)
}

// GetTags get all tags from a given catalogsService
// It is able to get tags from sub entities (like containers and services if needed)
func (s *CatalogServices) GetTags(c echo.Context) error {
	catalogService := c.Get("catalogService").(types.CatalogService)
	docktorAPI := c.Get("api").(*storage.Docktor)
	tagIds := catalogService.Tags

	tags, err := docktorAPI.Tags().FindAllByIDs(tagIds)
	if err != nil {
		log.WithError(err).WithField("catalogService", catalogService.ID).Error("Can't get tags of catalogService")
		return c.JSON(http.StatusInternalServerError, "Incorrect data. Contact your administrator")
	}
	return c.JSON(http.StatusOK, tags)
}
