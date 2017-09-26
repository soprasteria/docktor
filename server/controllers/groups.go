package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/soprasteria/docktor/server/models"
	"github.com/soprasteria/docktor/server/modules/users"
	"github.com/soprasteria/docktor/server/types"
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
		return c.String(http.StatusInternalServerError, "Error while retreiving all groups")
	}
	return c.JSON(http.StatusOK, groups)
}

//Save group into docktor
func (g *Groups) Save(c echo.Context) error {
	docktorAPI := c.Get("api").(*models.Docktor)
	var group types.Group
	err := c.Bind(&group)

	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error while binding group: %v", err))
	}

	// If the ID is empty, it's a creation, so generate an object ID
	if group.ID.Hex() == "" {
		group.ID = bson.NewObjectId()
	}

	// Filters the members by existing users
	// This way, group will autofix when user is deleted
	existingMembers := existingMembers(docktorAPI, group.Members)
	group.Members = existingMembers

	res, err := docktorAPI.Groups().Save(group)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error while saving group: %v", err))
	}
	return c.JSON(http.StatusOK, res)
}

// existingMembers return members filters by existing ones
// Checks wether the user actually exists in database
func existingMembers(docktorAPI *models.Docktor, members types.Members) types.Members {

	var existingMembers types.Members

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

//Delete group into docktor
func (g *Groups) Delete(c echo.Context) error {
	docktorAPI := c.Get("api").(*models.Docktor)
	id := c.Param("groupID")
	res, err := docktorAPI.Groups().Delete(bson.ObjectIdHex(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error while remove group: %v", err))
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
	withServices, _ := strconv.ParseBool(c.QueryParam("services"))     // Get all tags from a given daemon
	withcontainers, _ := strconv.ParseBool(c.QueryParam("containers")) // Get all tags from a given Users
	group := c.Get("group").(types.Group)
	docktorAPI := c.Get("api").(*models.Docktor)
	tagIds := group.Tags

	// Get also tags from container instances of group
	if withcontainers {
		for _, c := range group.Containers {
			tagIds = append(tagIds, c.Tags...)
		}
	}
	// Get also tags from the type of containers (= service)
	if withServices {
		var serviceIds []bson.ObjectId
		// Get services from containers
		for _, c := range group.Containers {
			serviceIds = append(serviceIds, c.ServiceID)
		}
		services, err := docktorAPI.Services().FindAllByIDs(serviceIds)
		if err != nil {
			log.WithError(err).WithField("group", group.ID).WithField("services.ids", serviceIds).Error("Can't get tags of service")
			return c.JSON(http.StatusInternalServerError, "Incorrect data. Contact your administrator")
		}
		// Get tags from services
		for _, s := range services {
			tagIds = append(tagIds, s.Tags...)
		}
	}

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
	for _, c := range group.Containers {
		daemonIds = append(daemonIds, c.DaemonID)
	}

	daemons, err := docktorAPI.Daemons().FindAllByIDs(daemonIds)
	if err != nil {
		log.WithError(err).WithField("group", group.ID).WithField("daemons.ids", daemonIds).Error("Can't get daemons of group")
		return c.JSON(http.StatusInternalServerError, "Incorrect data. Contact your administrator")
	}

	return c.JSON(http.StatusOK, daemons)
}

// GetServices get all services used on the group (service from containers)
func (g *Groups) GetServices(c echo.Context) error {
	group := c.Get("group").(types.Group)
	docktorAPI := c.Get("api").(*models.Docktor)

	var serviceIds []bson.ObjectId
	for _, c := range group.Containers {
		serviceIds = append(serviceIds, c.ServiceID)
	}

	services, err := docktorAPI.Services().FindAllByIDs(serviceIds)
	if err != nil {
		log.WithError(err).WithField("group", group.ID).WithField("service.ids", serviceIds).Error("Can't get services of group")
		return c.JSON(http.StatusInternalServerError, "Incorrect data. Contact your administrator")
	}
	return c.JSON(http.StatusOK, services)
}
