package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/soprasteria/docktor/server/models"
	"github.com/soprasteria/docktor/server/types"
	"gopkg.in/mgo.v2/bson"
)

// Tags contains all group handlers
type Tags struct {
}

//GetAll tags from docktor
func (s *Tags) GetAll(c echo.Context) error {
	docktorAPI := c.Get("api").(*models.Docktor)
	tags, err := docktorAPI.Tags().FindAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error while retreiving all tags")
	}
	return c.JSON(http.StatusOK, tags)
}

//Save or update tag into docktor
func (s *Tags) Save(c echo.Context) error {
	docktorAPI := c.Get("api").(*models.Docktor)
	id := c.Param("id")

	var tag types.Tag
	err := c.Bind(&tag)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Unable to parse the tag received from client: %v", err))
	}

	// Force ID in tag to be the one passed as parameter
	if id != "" {
		if !bson.IsObjectIdHex(id) {
			return c.String(http.StatusBadRequest, fmt.Sprintf("The ID %q is not a valid BSON id", id))
		}
		tag.ID = bson.ObjectIdHex(id)
	}

	res, err := docktorAPI.Tags().Save(tag)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("An error has occured while saving the tag: %v", err))
	}
	return c.JSON(http.StatusOK, res)
}

//Delete tag into docktor
func (s *Tags) Delete(c echo.Context) error {

	docktorAPI := c.Get("api").(*models.Docktor)
	id := c.Param("id")

	collections := []types.UseTags{
		docktorAPI.Daemons(),
		docktorAPI.Services(),
		// TODO : add others collections (users, groups and containers in groups)
	}

	// Remove tags from all collections containings tags
	// Don't fail the process even if one error occurs
	for _, c := range collections {
		changes, err := c.RemoveTag(bson.ObjectIdHex(id))
		if err != nil {
			log.WithError(err).WithField("tag", id).WithField("collection", c.GetCollectionName()).
				Error("Can't delete Removed tags of collection. Continuing anyway...")
		} else {
			log.WithField("tag", id).WithField("collection", c.GetCollectionName()).
				WithField("number_of_documents_updated", changes.Updated).
				Debug("Deleting tag : removed them from collection")
		}
	}

	res, err := docktorAPI.Tags().Delete(bson.ObjectIdHex(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error while deleting tag: %v", err))
	}
	return c.String(http.StatusOK, res.Hex())
}

// existingTags return tags filtered by existing ones
// Checks wether the tag actually exists in database
func existingTags(docktorAPI *models.Docktor, tagsIds []bson.ObjectId) []bson.ObjectId {

	existingTagsIDs := []bson.ObjectId{}

	// Get all real tags
	existingTags, _ := docktorAPI.Tags().FindAllByIDs(tagsIds)

	// Get their ids only
	for _, tag := range existingTags {
		existingTagsIDs = append(existingTagsIDs, tag.ID)
	}

	return existingTagsIDs
}
