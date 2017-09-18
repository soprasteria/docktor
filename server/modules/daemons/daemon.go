package daemons

import (
	"time"

	"github.com/soprasteria/docktor/server/adapters/cache"
	"github.com/soprasteria/docktor/server/types"
	"github.com/soprasteria/docktor/server/utils"
)

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
func GetInfo(daemon types.Daemon, cache cache.Cache, force bool) (*DaemonInfo, error) {

	key := daemon.ID.Hex()
	if !force {
		value, _ := cache.Get(key)
		if info, ok := value.(DaemonInfo); ok {
			return &info, nil
		}
	}

	api, err := utils.InitDocker(daemon)
	if err != nil {
		return nil, err
	}

	info := DaemonInfo{}

	dockerInfo, err := api.Docker.Info()
	if err != nil {
		info = DaemonInfo{Status: statusDOWN, NbImages: 0, NbContainers: 0, Message: err.Error()}
	} else {
		info = DaemonInfo{Status: statusUP, NbImages: dockerInfo.Images, NbContainers: dockerInfo.Containers}
	}

	go cache.Set(key, info, 5*time.Minute)
	return &info, nil
}
