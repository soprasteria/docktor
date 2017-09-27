package storage

import (
	"fmt"

	"github.com/soprasteria/docktor/server/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// TagsRepo is the repo for tags
type TagsRepo interface {
	// Save a tag into database
	Save(tag types.Tag) (types.Tag, error)
	// Delete a tag in database
	Delete(id bson.ObjectId) (bson.ObjectId, error)
	// FindByID get the tag by its id
	FindByID(id string) (types.Tag, error)
	// FindByIDBson get the tag by its id
	FindByIDBson(id bson.ObjectId) (types.Tag, error)
	// Find get the first tag with a given name and category
	Find(name string, category string) (types.Tag, error)
	// FindAll get all tags
	FindAll() ([]types.Tag, error)
	// FindAllByIDs get all tags with id
	FindAllByIDs([]bson.ObjectId) ([]types.Tag, error)
	// Drop drops the content of the collection
	Drop() error
	// GetCollectionName returns the name of the collection
	GetCollectionName() string
}

// DefaultTagsRepo is the repository for tags
type DefaultTagsRepo struct {
	coll *mgo.Collection
}

// NewTagsRepo instantiate new RepoDaemons
func NewTagsRepo(coll *mgo.Collection) TagsRepo {
	return &DefaultTagsRepo{coll: coll}
}

// CreateIndexes creates Index
func (r *DefaultTagsRepo) CreateIndexes() error {
	return r.coll.EnsureIndex(mgo.Index{
		Key:    []string{"category", "name"},
		Unique: true,
		Name:   "tag_cat_name_unique",
	})
}

// GetCollectionName gets the name of the collection
func (r *DefaultTagsRepo) GetCollectionName() string {
	return r.coll.FullName
}

// Save or create a tag into a database
// Two different tags can not have the same name and category (unicity is checked on slugified name and category)
func (r *DefaultTagsRepo) Save(tag types.Tag) (types.Tag, error) {

	_, err := r.coll.UpsertId(tag.ID, bson.M{"$set": tag})
	if mgo.IsDup(err) {
		return tag, fmt.Errorf("Another tag exists with category %v and name %v", tag.Category.GetRaw(), tag.Name.GetRaw())
	}

	return tag, err
}

// Delete a tag in database
func (r *DefaultTagsRepo) Delete(id bson.ObjectId) (bson.ObjectId, error) {
	err := r.coll.RemoveId(id)
	return id, err
}

// FindByID get the tag by its id
func (r *DefaultTagsRepo) FindByID(id string) (types.Tag, error) {
	result := types.Tag{}
	err := r.coll.FindId(bson.ObjectIdHex(id)).One(&result)
	return result, err
}

// FindByIDBson get the tag by its id (as a bson object)
func (r *DefaultTagsRepo) FindByIDBson(id bson.ObjectId) (types.Tag, error) {
	result := types.Tag{}
	err := r.coll.FindId(id).One(&result)
	return result, err
}

// findBySlug get tag identified by its name slugified
func (r *DefaultTagsRepo) findBySlug(slugName string, slugCategory string) (types.Tag, error) {
	result := types.Tag{}
	err := r.coll.Find(bson.M{"name.slug": slugName, "category.slug": slugCategory}).One(&result)
	return result, err
}

// Find get the first tag with a given name
func (r *DefaultTagsRepo) Find(name string, category string) (types.Tag, error) {
	tagName := types.NewTagName(name)
	tagCategory := types.NewTagCategory(category)
	return r.findBySlug(tagName.GetSlug(), tagCategory.GetSlug())
}

// FindAll get all tags
func (r *DefaultTagsRepo) FindAll() ([]types.Tag, error) {
	results := []types.Tag{}
	err := r.coll.Find(bson.M{}).All(&results)
	return results, err
}

// FindAllByIDs get all tags with id
func (r *DefaultTagsRepo) FindAllByIDs(ids []bson.ObjectId) ([]types.Tag, error) {
	results := []types.Tag{}
	err := r.coll.Find(
		bson.M{"_id": &bson.M{"$in": ids}},
	).All(&results)
	return results, err
}

// Drop drops the content of the collection
func (r *DefaultTagsRepo) Drop() error {
	return r.coll.DropCollection()
}
