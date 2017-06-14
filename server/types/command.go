package types

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Command for images
type Command struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string        `bson:"name" json:"name"`
	Exec       string        `bson:"exec" json:"exec"`             // Effective command to execute
	Parameters []string      `bson:"parameters" json:"parameters"` // Parameters that can be filled in by user at runtime
	Roles      MemberRole    `bson:"role" json:"role"`             // Only members with one of these roles can execute the command
	Created    time.Time     `bson:"created" json:"created"`
	Updated    time.Time     `bson:"updated" json:"updated"`
}

// Commands is a slice of Command
type Commands []Command
