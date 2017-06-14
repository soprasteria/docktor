package types

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Template is an archetype to bootstrap a set of services
// These services are not bound together, meaning a service can later be deleted once instanciated
type Template struct {
	ID          bson.ObjectId   `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string          `bson:"name" json:"name"`
	Description string          `bson:"description" json:"description"`
	Services    []Service       `bson:"services" json:"services"`
	Tags        []bson.ObjectId `bson:"tags" json:"tags"`
	Created     time.Time       `bson:"created" json:"created"`
	Updated     time.Time       `bson:"updated" json:"updated"`
}
