package daemons

import (
	"time"

	"github.com/soprasteria/docktor/server/types"
	"github.com/soprasteria/docktor/server/utils"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/redis.v3"
)

// DaemonRest is a simplified daemon. It is meant to be fetched by user that are not admins
// This kind of simple user should not have access to protected fields like certificates and so on.
type DaemonRest struct {
	ID          bson.ObjectId   `bson:"_id,omitempty" json:"id,omitempty"`
	Active      bool            `bson:"active" json:"active"`
	Name        string          `bson:"name" json:"name" validate:"required"`
	Description string          `bson:"description,omitempty" json:"description,omitempty"`
	Site        bson.ObjectId   `bson:"site" json:"site" validate:"required"`
	Variables   types.Variables `bson:"variables" json:"variables"`
	Volumes     types.Volumes   `bson:"volumes" json:"volumes"`
	Tags        []bson.ObjectId `bson:"tags" json:"tags"`
	Created     time.Time       `bson:"created" json:"created"` // Fields that will be populated automatically by server
	Updated     time.Time       `bson:"updated" json:"updated"` // Fields that will be populated automatically by server
}

// GetDaemonRest returns a Docktor daemon, amputed of sensible data
func GetDaemonRest(d types.Daemon) DaemonRest {
	return DaemonRest{
		ID:          d.ID,
		Active:      d.Active,
		Name:        d.Name,
		Description: d.Description,
		Variables:   d.Variables,
		Volumes:     d.Volumes,
		Tags:        d.Tags,
		Created:     d.Created,
		Updated:     d.Updated,
	}
}

// GetDaemonsRest returns the slice of Docktor daemon, amputed of sensible data
func GetDaemonsRest(daemons []types.Daemon) []DaemonRest {
	daemonsRest := []DaemonRest{}
	for _, v := range daemons {
		daemonsRest = append(daemonsRest, GetDaemonRest(v))
	}
	return daemonsRest
}

// DaemonInfo struct
type DaemonInfo struct {
	Status       string `json:"status"`
	NbImages     int    `json:"nbImages"`
	NbContainers int    `json:"nbContainers"`
	Message      string `json:"message,omitempty"`
}

const (
	statusUP   string = "UP"
	statusDOWN string = "DOWN"
)

// GetInfo : retrieving the docker daemon status using redis cache
func GetInfo(daemon types.Daemon, client *redis.Client, force bool) (*DaemonInfo, error) {
	info := &DaemonInfo{}
	key := daemon.ID.Hex()
	if !force {
		err := utils.GetFromRedis(client, key, info)
		if err == nil {
			return info, nil
		}
	}

	api, err := utils.InitDocker(daemon)
	if err != nil {
		return nil, err
	}

	dockerInfo, err := api.Docker.Info()
	if err != nil {
		info = &DaemonInfo{Status: statusDOWN, NbImages: 0, NbContainers: 0, Message: err.Error()}
		go utils.SetIntoRedis(client, key, info, 5*time.Minute)
		return info, nil
	}

	info = &DaemonInfo{Status: statusUP, NbImages: dockerInfo.Images, NbContainers: dockerInfo.Containers}
	go utils.SetIntoRedis(client, key, info, 5*time.Minute)
	return info, nil
}
