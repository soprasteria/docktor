package storage

import (
	"fmt"

	"github.com/soprasteria/docktor/server/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CatalogServicesRepo is the repo for CatalogService
type CatalogServicesRepo interface {
	//===========
	// CatalogServices
	//===========

	// Drop drops the content of the collection
	Drop() error

	// Save a CatalogService into database
	Save(catalogService types.CatalogService) (types.CatalogService, error)

	// Delete a CatalogService in database
	Delete(id bson.ObjectId) (bson.ObjectId, error)

	// FindByID get the CatalogService by its id
	FindByID(id string) (types.CatalogService, error)

	// FindByIDBson get the CatalogService by its id
	FindByIDBson(id bson.ObjectId) (types.CatalogService, error)

	// Find get the first CatalogService with a given name
	Find(name string) (types.CatalogService, error)

	// FindAll get all CatalogServices
	FindAll() ([]types.CatalogService, error)

	// FindAllByName get all CatalogServices by the give name
	FindAllByName(name string) ([]types.CatalogService, error)

	// FindAllByIDs get all CatalogServices from their ids
	FindAllByIDs(ids []bson.ObjectId) ([]types.CatalogService, error)

	// GetCollectionName returns the name of the collection
	GetCollectionName() string

	// RemoveTag
	RemoveTag(id bson.ObjectId) (*mgo.ChangeInfo, error)
}

// DefaultCatalogServicesRepo is the repository for catalogServices
type DefaultCatalogServicesRepo struct {
	coll *mgo.Collection
}

// NewCatalogServicesRepo instantiate new CatalogServicesRepo
func NewCatalogServicesRepo(coll *mgo.Collection) CatalogServicesRepo {
	return &DefaultCatalogServicesRepo{coll: coll}
}

// GetCollectionName gets the name of the collection
func (r *DefaultCatalogServicesRepo) GetCollectionName() string {
	return r.coll.FullName
}

// CreateIndexes creates Index
func (r *DefaultCatalogServicesRepo) CreateIndexes() error {
	return r.coll.EnsureIndex(mgo.Index{
		Key:    []string{"title"},
		Unique: true,
		Name:   "catalogService_title_unique",
	})
}

// Drop drops the content of the collection
func (r *DefaultCatalogServicesRepo) Drop() error {
	return r.coll.DropCollection()
}

// Save a catalogService into a database
func (r *DefaultCatalogServicesRepo) Save(catalogService types.CatalogService) (types.CatalogService, error) {
	newCatalogService := types.NewCatalogService(catalogService)
	_, err := r.coll.UpsertId(catalogService.ID, bson.M{"$set": types.NewCatalogService(catalogService)})
	if mgo.IsDup(err) {
		return catalogService, fmt.Errorf("Another catalogService exists with name '%v'", catalogService.Name)
	}
	return newCatalogService, err
}

// Delete a catalogService in database
func (r *DefaultCatalogServicesRepo) Delete(id bson.ObjectId) (bson.ObjectId, error) {
	err := r.coll.RemoveId(id)
	return id, err
}

// Find get the first catalogService with a given name
func (r *DefaultCatalogServicesRepo) Find(name string) (types.CatalogService, error) {
	result := types.CatalogService{}
	err := r.coll.Find(bson.M{"title": name}).One(&result)
	return result, err
}

// FindByID get the catalogService by its id
func (r *DefaultCatalogServicesRepo) FindByID(id string) (types.CatalogService, error) {
	result := types.CatalogService{}
	err := r.coll.FindId(bson.ObjectIdHex(id)).One(&result)
	return result, err
}

// FindByIDBson get the catalogService by its id (as a bson object)
func (r *DefaultCatalogServicesRepo) FindByIDBson(id bson.ObjectId) (types.CatalogService, error) {
	result := types.CatalogService{}
	err := r.coll.FindId(id).One(&result)
	return result, err
}

// FindAll get all catalogServices
func (r *DefaultCatalogServicesRepo) FindAll() ([]types.CatalogService, error) {
	results := []types.CatalogService{}
	err := r.coll.Find(bson.M{}).All(&results)
	return results, err
}

// FindAllByIDs get all catalogServices from thei ids
func (r *DefaultCatalogServicesRepo) FindAllByIDs(ids []bson.ObjectId) ([]types.CatalogService, error) {
	results := []types.CatalogService{}
	err := r.coll.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&results)
	return results, err
}

// FindAllByName get all catalogServices by the give name
func (r *DefaultCatalogServicesRepo) FindAllByName(name string) ([]types.CatalogService, error) {
	results := []types.CatalogService{}
	err := r.coll.Find(bson.M{"title": name}).All(&results)
	return results, err
}

// RemoveTag removes given tag from all users
func (r *DefaultCatalogServicesRepo) RemoveTag(id bson.ObjectId) (*mgo.ChangeInfo, error) {
	return r.coll.UpdateAll(
		bson.M{"tags": bson.M{"$in": []bson.ObjectId{id}}},
		bson.M{"$pull": bson.M{"tags": id}},
	)
}
