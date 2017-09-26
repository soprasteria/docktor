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

// Tags contains all group handlers
type Tags struct {
}

// TagAlreadyExistErrMessage is an error message when a tag alread exists
const TagAlreadyExistErrMessage string = "Tag %q already exists in category %q"

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
	var tag types.Tag
	err := c.Bind(&tag)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Unable to parse the tag received from client: %v", err))
	}

	tagByNameAndCat, err := docktorAPI.Tags().Find(tag.Name.GetRaw(), tag.Category.GetRaw())

	// Update fields
	id := c.Param("id")
	if tag.ID.Hex() == "" && id == "" {
		// New tag to create
		if err == nil {
			// Tag cannot created when we found another tag with same info
			log.WithField("newtag_name", tag.Name).WithField("newtag_category", tag.Category).
				WithField("existingtag_name", tagByNameAndCat.Name).WithField("existingtag_category", tagByNameAndCat.Category).
				Warning("Can't create tag because it already exists...")
			return fmt.Errorf(
				TagAlreadyExistErrMessage,
				tag.Name.GetRaw(), tag.Category.GetRaw(),
			)
		}
		tag.ID = bson.NewObjectId()
		tag.Created = time.Now()
	} else {
		// Existing tag, search for it and update read-only fields
		tag.ID = bson.ObjectIdHex(id)
		t, errr := docktorAPI.Tags().FindByIDBson(tag.ID)
		if errr != nil {
			if errr == mgo.ErrNotFound {
				return c.String(http.StatusBadRequest, fmt.Sprint("Tag does not exist"))
			}
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to find tag. Retry later : %s", errr))
		}
		tag.Created = t.Created
	}
	tag.Updated = time.Now()

	// Validate fields from validator tags for common types
	if err = c.Validate(tag); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Some fields of tag are not valid: %v", err))
	}
	res, err := docktorAPI.Tags().Save(tag)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("An error has occurred while saving the tag: %v", err))
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
		docktorAPI.Groups(),
		docktorAPI.Users(),
		// TODO : add others collections (containers in groups)
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
