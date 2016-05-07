package controllers

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	api "github.com/soprasteria/godocktor-api"
	"github.com/soprasteria/godocktor-api/types"
	"gopkg.in/mgo.v2/bson"
)

// GroupsController contains all groups handlers
type GroupsController struct {
}

//GetAllGroups from docktor
func (gc *GroupsController) GetAllGroups(c echo.Context) error {
	docktorAPI := c.Get("api").(*api.Docktor)
	groups, err := docktorAPI.Groups().FindAll()
	glog.Info("test")
	if err != nil {
		return c.String(500, "Error while retreiving all groups")
	}
	return c.JSON(200, groups)
}

//SaveGroup into docktor
func (gc *GroupsController) SaveGroup(c echo.Context) error {
	docktorAPI := c.Get("api").(*api.Docktor)
	var group types.Group
	err := c.Bind(&group)

	if err != nil {
		return c.String(400, fmt.Sprintf("Error while binding group: %v", err))
	}
	res, err := docktorAPI.Groups().Save(group)
	if err != nil {
		return c.String(500, fmt.Sprintf("Error while saving group: %v", err))
	}
	return c.JSON(200, res)
}

//DeleteGroup into docktor
func (gc *GroupsController) DeleteGroup(c echo.Context) error {
	docktorAPI := c.Get("api").(*api.Docktor)
	id := c.Param("id")
	res, err := docktorAPI.Groups().Delete(bson.ObjectIdHex(id))
	if err != nil {
		return c.String(500, fmt.Sprintf("Error while remove group: %v", err))
	}
	return c.JSON(200, res)
}
