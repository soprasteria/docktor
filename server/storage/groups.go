package storage

import (
	"fmt"

	"github.com/soprasteria/docktor/server/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GroupsRepo is the repo for groups
type GroupsRepo interface {
	//===========
	// Groups
	//===========

	// Drop drops the content of the collection
	Drop() error
	// Save a group into database
	Save(group types.Group) (types.Group, error)
	// Delete a group in database
	Delete(id bson.ObjectId) (bson.ObjectId, error)
	// FindByID get the group by its id
	FindByID(id string) (types.Group, error)
	// FindByIDBson get the group by its id
	FindByIDBson(id bson.ObjectId) (types.Group, error)
	// Find get the first group with a given name
	Find(name string) (types.Group, error)
	// FindAll get all groups
	FindAll() ([]types.Group, error)
	// FindAllByName get all groups by the give name
	FindAllByName(name string) ([]types.Group, error)
	// FindAllByIDs get all groups from thei ids
	FindAllByIDs(ids []bson.ObjectId) ([]types.Group, error)
	// FindAllByRegex get all groups by the regex name
	FindAllByRegex(nameRegex string) ([]types.Group, error)
	// FindAllWithContainers get all groups that contains a list of containers
	FindAllWithContainers(groupNameRegex string, containersID []string) ([]types.Group, error)
	// RemoveMember remove a member from all groups
	RemoveMember(userID bson.ObjectId) (*mgo.ChangeInfo, error)
	// GetCollectionName returns the name of the collection
	GetCollectionName() string
	// CreateIndexes creates Index
	CreateIndexes() error
	// RemoveTag
	RemoveTag(id bson.ObjectId) (*mgo.ChangeInfo, error)
	// RemoveFileSystem removes filesystems with given daemon from all groups
	RemoveFileSystem(daemonID bson.ObjectId) (*mgo.ChangeInfo, error)
}

// DefaultGroupsRepo is the repository for groups
type DefaultGroupsRepo struct {
	coll *mgo.Collection
}

// NewGroupsRepo instantiate new GroupsRepo
func NewGroupsRepo(coll *mgo.Collection) GroupsRepo {
	return &DefaultGroupsRepo{coll: coll}
}

// GetCollectionName gets the name of the collection
func (r *DefaultGroupsRepo) GetCollectionName() string {
	return r.coll.FullName
}

// CreateIndexes creates Index
func (r *DefaultGroupsRepo) CreateIndexes() error {
	return r.coll.EnsureIndex(mgo.Index{
		Key:    []string{"title"},
		Unique: true,
		Name:   "group_title_unique",
	})
}

// Drop drops the content of the collection
func (r *DefaultGroupsRepo) Drop() error {
	return r.coll.DropCollection()
}

// Save a group into a database
func (r *DefaultGroupsRepo) Save(group types.Group) (types.Group, error) {
	newGroup := types.NewGroup(group)
	_, err := r.coll.UpsertId(group.ID, bson.M{"$set": types.NewGroup(group)})
	if mgo.IsDup(err) {
		return group, fmt.Errorf("Another group exists with title '%v'", group.Title)
	}
	return newGroup, err
}

// Delete a group in database
func (r *DefaultGroupsRepo) Delete(id bson.ObjectId) (bson.ObjectId, error) {
	err := r.coll.RemoveId(id)
	return id, err
}

// Find get the first group with a given name
func (r *DefaultGroupsRepo) Find(name string) (types.Group, error) {
	result := types.Group{}
	err := r.coll.Find(bson.M{"title": name}).One(&result)
	return result, err
}

// FindByID get the group by its id
func (r *DefaultGroupsRepo) FindByID(id string) (types.Group, error) {
	result := types.Group{}
	err := r.coll.FindId(bson.ObjectIdHex(id)).One(&result)
	return result, err
}

// FindByIDBson get the group by its id (as a bson object)
func (r *DefaultGroupsRepo) FindByIDBson(id bson.ObjectId) (types.Group, error) {
	result := types.Group{}
	err := r.coll.FindId(id).One(&result)
	return result, err
}

// FindAll get all groups
func (r *DefaultGroupsRepo) FindAll() ([]types.Group, error) {
	results := []types.Group{}
	err := r.coll.Find(bson.M{}).All(&results)
	return results, err
}

// FindAllByIDs get all groups from thei ids
func (r *DefaultGroupsRepo) FindAllByIDs(ids []bson.ObjectId) ([]types.Group, error) {
	results := []types.Group{}
	err := r.coll.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&results)
	return results, err
}

// FindAllByName get all groups by the give name
func (r *DefaultGroupsRepo) FindAllByName(name string) ([]types.Group, error) {
	results := []types.Group{}
	err := r.coll.Find(bson.M{"title": name}).All(&results)
	return results, err
}

// FindAllByRegex get all groups by the regex name
func (r *DefaultGroupsRepo) FindAllByRegex(nameRegex string) ([]types.Group, error) {
	results := []types.Group{}
	err := r.coll.Find(bson.M{"title": &bson.RegEx{Pattern: nameRegex}}).All(&results)
	return results, err
}

// FindAllWithContainers get all groups that contains a list of containers
func (r *DefaultGroupsRepo) FindAllWithContainers(groupNameRegex string, containersID []string) ([]types.Group, error) {
	results := []types.Group{}
	err := r.coll.Find(
		bson.M{
			"title":                  &bson.RegEx{Pattern: groupNameRegex},
			"containers.containerId": &bson.M{"$in": containersID},
		}).All(&results)

	return results, err
}

// RemoveMember remove a member from all groups
func (r *DefaultGroupsRepo) RemoveMember(userID bson.ObjectId) (*mgo.ChangeInfo, error) {
	return r.coll.UpdateAll(
		bson.M{},
		bson.M{"$pull": bson.M{"members": bson.M{"user": userID}}},
	)
}

// RemoveTag removes given tag from all users
func (r *DefaultGroupsRepo) RemoveTag(id bson.ObjectId) (*mgo.ChangeInfo, error) {
	return r.coll.UpdateAll(
		bson.M{"tags": bson.M{"$in": []bson.ObjectId{id}}},
		bson.M{"$pull": bson.M{"tags": id}},
	)
}

// RemoveFileSystem removes filesystems with given daemon from all groups
func (r *DefaultGroupsRepo) RemoveFileSystem(daemonID bson.ObjectId) (*mgo.ChangeInfo, error) {
	return r.coll.UpdateAll(
		bson.M{"filesystems.daemon": bson.M{"$in": []bson.ObjectId{daemonID}}},
		bson.M{"$pull": bson.M{"filesystems": bson.M{"daemon": daemonID}}},
	)
}
