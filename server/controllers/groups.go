package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/soprasteria/docktor/server/models"
	"github.com/soprasteria/docktor/server/modules/daemons"
	"github.com/soprasteria/docktor/server/modules/users"
	"github.com/soprasteria/docktor/server/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Groups contains all groups handlers
type Groups struct {
}

//GetAll groups from docktor
func (g *Groups) GetAll(c echo.Context) error {
	docktorAPI := c.Get("api").(*models.Docktor)
	groups, err := docktorAPI.Groups().FindAll()
	if err != nil {
		log.WithError(err).Error("Unable to get all groups")
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to get all daemons because of technical error: %v. Retry later.", err))
	}
	return c.JSON(http.StatusOK, groups)
}

//Save group into docktor
func (g *Groups) Save(c echo.Context) error {
	docktorAPI := c.Get("api").(*models.Docktor)
	var group types.Group
	if err := c.Bind(&group); err != nil {
		log.WithError(err).Error("Unable to bind group to save")
		return c.String(http.StatusBadRequest, "Unable to parse group received from client")
	}

	// If the ID is empty, it's a creation, so generate an object ID
	id := c.Param("groupID")
	if id == "" {
		// New group to create
		group.ID = bson.NewObjectId()
		group.Created = time.Now()
	} else {
		// Existing group
		group.ID = bson.ObjectIdHex(id)
		g, err := docktorAPI.Groups().FindByIDBson(group.ID)
		if err != nil {
			if err == mgo.ErrNotFound {
				log.WithError(err).Warnf("Tried to save a group that does not exist: %v", group.ID)
				return c.String(http.StatusBadRequest, "Group does not exist")
			}
			log.WithError(err).Errorf("Unable to find group because of unexpected error : %v", group.ID)
			return c.String(http.StatusInternalServerError, "Unable to find group because of technical error. Retry later.")
		}
		group.Created = g.Created
	}

	// Validate fields from validator tags for common types
	if err := c.Validate(group); err != nil {
		log.WithError(err).Errorf("Unable to save group %v because some fields are not valid", group.ID)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Some fields of group are not valid: %v", err))
	}

	// Validate fields that cannot be validated by validator engine
	if err := group.Validate(); err != nil {
		log.WithError(err).Errorf("Unable to save group %v because some fields are not valid", group.ID)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Some fields of group are not valid: %v", err))
	}

	// Keep only existing and remove duplicates of external collections
	// Used to clean the group of old data when saving
	group.Members = existingMembers(docktorAPI, group.Members)
	group.Tags = existingTags(docktorAPI, group.Tags)
	group.FileSystems = existingFileSystems(docktorAPI, group.FileSystems)
	group.Updated = time.Now()

	res, err := docktorAPI.Groups().Save(group)
	if err != nil {
		log.WithError(err).Errorf("Unexpected error when saving group %v", group.ID)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to save group because of technical error: %v. Retry later.", err))
	}
	return c.JSON(http.StatusOK, res)
}

// existingMembers return members filters by existing ones
// Checks wether the user actually exists in database
func existingMembers(docktorAPI *models.Docktor, members types.Members) types.Members {

	existingMembers := types.Members{}

	// Get all real users from members.
	existingUsers, _ := docktorAPI.Users().FindAllByIDs(members.GetUsers())

	// x*x Cardinality because no need to overoptimize with maps
	// as we will not have millions of members in a group
	for _, user := range existingUsers {
		for _, member := range members {
			if user.ID == member.User {
				existingMembers = append(existingMembers, member)
			}
		}
	}

	return existingMembers
}

// existingGroups return groups filtered by existing ones
// Checks wether the group actually exists in database
func existingGroups(docktorAPI *models.Docktor, groupsIDs []bson.ObjectId) []bson.ObjectId {

	existingGroupIDs := []bson.ObjectId{}

	// Get all real groups
	existingGroups, _ := docktorAPI.Groups().FindAllByIDs(groupsIDs)

	// Get their ids only
	for _, g := range existingGroups {
		existingGroupIDs = append(existingGroupIDs, g.ID)
	}

	return existingGroupIDs
}

// existingFileSystems return filesystems filtered by existing ones
// Checks wether the filesystem daemon actually exists in database
func existingFileSystems(docktorAPI *models.Docktor, fileSystems types.FileSystems) types.FileSystems {

	existingFileSystems := types.FileSystems{}
	daemonIDs := []bson.ObjectId{}
	for _, r := range fileSystems {
		daemonIDs = append(daemonIDs, r.Daemon)
	}

	// Get all real tags
	existingDaemons, _ := docktorAPI.Daemons().FindAllByIDs(daemonIDs)

	// Get existing filesystems only
	for _, daemon := range existingDaemons {
		for _, fs := range fileSystems {
			if daemon.ID == fs.Daemon {
				existingFileSystems = append(existingFileSystems, fs)
			}
		}
	}

	return existingFileSystems
}

//Delete group into docktor
func (g *Groups) Delete(c echo.Context) error {
	docktorAPI := c.Get("api").(*models.Docktor)
	id := c.Param("groupID")

	// TODO : Don't delete group if at least a service is deployed somewhere.

	res, err := docktorAPI.Groups().Delete(bson.ObjectIdHex(id))
	if err != nil {
		log.WithError(err).Errorf("Unexpected error when deleting group %v", id)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to delete group because of technical error: %v. Retry later.", err))
	}
	return c.String(http.StatusOK, res.Hex())
}

//Get group from docktor
func (g *Groups) Get(c echo.Context) error {
	group := c.Get("group").(types.Group)
	return c.JSON(http.StatusOK, group)
}

// GetTags get all tags from a given group
// It is able to get get tags from sub entities (like containers and services if needed)
func (g *Groups) GetTags(c echo.Context) error {
	// withServices, _ := strconv.ParseBool(c.QueryParam("services"))     // Get all tags from a given daemon
	// withcontainers, _ := strconv.ParseBool(c.QueryParam("containers")) // Get all tags from a given Users
	group := c.Get("group").(types.Group)
	docktorAPI := c.Get("api").(*models.Docktor)
	tagIds := group.Tags

	// TODO : enable it when containers and services are used again.
	// Get also tags from container instances of group
	// if withcontainers {
	// 	for _, c := range group.Containers {
	// 		tagIds = append(tagIds, c.Tags...)
	// 	}
	// }
	// // Get also tags from the type of containers (= service)
	// if withServices {
	// 	var serviceIds []bson.ObjectId
	// 	// Get services from containers
	// 	for _, c := range group.Containers {
	// 		serviceIds = append(serviceIds, c.ServiceID)
	// 	}
	// 	services, err := docktorAPI.Services().FindAllByIDs(serviceIds)
	// 	if err != nil {
	// 		log.WithError(err).WithField("group", group.ID).WithField("services.ids", serviceIds).Error("Can't get tags of service")
	// 		return c.JSON(http.StatusInternalServerError, "Incorrect data. Contact your administrator")
	// 	}
	// 	// Get tags from services
	// 	for _, s := range services {
	// 		tagIds = append(tagIds, s.Tags...)
	// 	}
	// }

	tags, err := docktorAPI.Tags().FindAllByIDs(tagIds)
	if err != nil {
		log.WithError(err).WithField("group", group.ID).Error("Can't get tags of group")
		return c.JSON(http.StatusInternalServerError, "Incorrect data. Contact your administrator")
	}

	return c.JSON(http.StatusOK, tags)
}

// GetMembers get all users who are members of the group
func (g *Groups) GetMembers(c echo.Context) error {
	group := c.Get("group").(types.Group)
	docktorAPI := c.Get("api").(*models.Docktor)

	ur := users.Rest{Docktor: docktorAPI}
	users, err := ur.GetUsersFromIds(group.Members.GetUsers())

	if err != nil {
		log.WithError(err).WithField("group", group.ID).Error("Can't get members of group")
		return c.JSON(http.StatusInternalServerError, "Incorrect data. Contact your administrator")
	}

	return c.JSON(http.StatusOK, users)
}

// GetDaemons get all daemons used on the group (filesystem and containers)
func (g *Groups) GetDaemons(c echo.Context) error {

	group := c.Get("group").(types.Group)
	docktorAPI := c.Get("api").(*models.Docktor)

	var daemonIds []bson.ObjectId

	for _, fs := range group.FileSystems {
		daemonIds = append(daemonIds, fs.Daemon)
	}

	// TODO : enable it when containers and services are used again.
	// for _, c := range group.Containers {
	// 	daemonIds = append(daemonIds, c.DaemonID)
	// }

	ds, err := docktorAPI.Daemons().FindAllByIDs(daemonIds)
	if err != nil {
		log.WithError(err).WithField("group", group.ID).WithField("daemons.ids", daemonIds).Error("Can't get daemons of group")
		return c.JSON(http.StatusInternalServerError, "Incorrect data. Contact your administrator")
	}

	return c.JSON(http.StatusOK, daemons.GetDaemonsRest(ds))
}

// GetServices get all services used on the group (service from containers)
func (g *Groups) GetServices(c echo.Context) error {
	group := c.Get("group").(types.Group)
	docktorAPI := c.Get("api").(*models.Docktor)

	serviceIds := []bson.ObjectId{}

	// TODO : enable it when containers and services are used again.
	// for _, c := range group.Containers {
	// 	serviceIds = append(serviceIds, c.ServiceID)
	// }

	services, err := docktorAPI.Services().FindAllByIDs(serviceIds)
	if err != nil {
		log.WithError(err).WithField("group", group.ID).WithField("service.ids", serviceIds).Error("Can't get services of group")
		return c.JSON(http.StatusInternalServerError, "Incorrect data. Contact your administrator")
	}
	return c.JSON(http.StatusOK, services)
}
