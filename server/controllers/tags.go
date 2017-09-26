package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/soprasteria/docktor/server/models"
	"github.com/soprasteria/docktor/server/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// TagAlreadyExistErrMessage is an error message when a tag alread exists
const TagAlreadyExistErrMessage string = "Tag %q already exists in category %q"

// Tags contains all group handlers
type Tags struct {
}

//GetAll tags from docktor
func (s *Tags) GetAll(c echo.Context) error {
	docktorAPI := c.Get("api").(*models.Docktor)
	tags, err := docktorAPI.Tags().FindAll()
	if err != nil {
		log.WithError(err).Error("Unable to get all tags")
		return c.String(http.StatusInternalServerError, "Unable to get all tags because of technical error. Retry later")
	}
	return c.JSON(http.StatusOK, tags)
}

//Save or update tag into docktor
func (s *Tags) Save(c echo.Context) error {
	docktorAPI := c.Get("api").(*models.Docktor)

	// Unserialize the tag
	var tagToSave types.Tag
	if err := c.Bind(&tagToSave); err != nil {
		log.WithError(err).Error("Unable to bind tag to save")
		return c.String(http.StatusBadRequest, "Unable to parse tag received from client")
	}

	id := c.Param("tagID")
	var savedTag types.Tag
	if id == "" {
		// Tag to create
		savedTag.ID = bson.NewObjectId()
		savedTag.Created = time.Now()
	} else {
		// Tag to update
		savedTag.ID = bson.ObjectIdHex(id)
		existingTag, err := docktorAPI.Tags().FindByIDBson(savedTag.ID)
		if err != nil {
			if err == mgo.ErrNotFound {
				log.WithError(err).Warnf("Tried to save a tag that does not exist: %v", savedTag.ID)
				return c.String(http.StatusBadRequest, "Tag does not exist")
			}
			log.WithError(err).Errorf("Unable to find tag because of unexpected error : %v", savedTag.ID)
			return c.String(http.StatusInternalServerError, "Unable to find tag because of technical error. Retry later.")
		}
		savedTag.Created = existingTag.Created
	}

	// Set values to writable data
	savedTag.Name = types.NewTagName(tagToSave.Name.GetRaw())
	savedTag.Category = types.NewTagCategory(tagToSave.Category.GetRaw())
	savedTag.Updated = time.Now()
	savedTag.UsageRights = tagToSave.UsageRights

	// Set default values
	if savedTag.UsageRights == "" {
		savedTag.UsageRights = types.AdminRole
	}

	// Validate fields
	if err := c.Validate(savedTag); err != nil {
		log.WithError(err).Errorf("Unable to save tag %v because some fields are not valid", savedTag.ID)
		return c.String(http.StatusBadRequest, "Category, name and usage rights are required")
	}
	if !savedTag.UsageRights.IsValid() {
		err := fmt.Errorf("Expected userRights to be 'admin' or 'user', obtained '%v'", savedTag.UsageRights)
		log.WithError(err).Errorf("Unable to save tag %v because usage rights are not valid", savedTag.ID)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Unable to save tag because usage rights are not valid: %v", err))
	}

	// Saving to database
	res, err := docktorAPI.Tags().Save(savedTag)
	if err != nil {
		log.WithError(err).Errorf("Unable to save tag %v because of technical error", savedTag.ID)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to save tag because of technical error: %v.", err))
	}
	return c.JSON(http.StatusOK, res)
}

//Delete tag into docktor
func (s *Tags) Delete(c echo.Context) error {

	docktorAPI := c.Get("api").(*models.Docktor)
	id := c.Param("tagID")

	collections := []types.UseTags{
		docktorAPI.Daemons(),
		docktorAPI.Services(),
		docktorAPI.Users(),
		// TODO : add others collections (groups, services ...)
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
		log.WithError(err).Errorf("Unable to delete tag %v because of database error", id)
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
