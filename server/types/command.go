package types

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Command for images
type Command struct {
	ID            bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name          string        `bson:"name" json:"name"`
	Exec          string        `bson:"exec" json:"exec"`                   // Effective command to execute
	AuthorizeArgs bool          `bson:"authorizeArgs" json:"authorizeArgs"` // When true, arguments can be added at runtime by user
	DefaultArgs   []string      `bson:"defaultArgs" json:"defaultArgs"`     // When set and AuthorizeArgs is true, defines the list of arguments that used by user at runtime.
	Roles         MemberRole    `bson:"role" json:"role"`                   // Only members with one of these roles can execute the command
	Created       time.Time     `bson:"created" json:"created"`
	Updated       time.Time     `bson:"updated" json:"updated"`
}

// Commands is a slice of Command
type Commands []Command
