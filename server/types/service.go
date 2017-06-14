package types

import (
	// "errors"
	// "fmt"
	// "strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Service defines a CDK service in a catalog, i.e. a service that can be deployed to a given machine
// A service contains many versions. Each version contains 1 or more containers
// Containers of a service are bound together, meaning if the service is stopped, all containers are stopped
// Jobs (checkhealth) or commands can be executed on containers
type Service struct {
	ID       bson.ObjectId    `bson:"_id,omitempty" json:"id,omitempty"`
	Created  time.Time        `bson:"created" json:"created"`
	Updated  time.Time        `bson:"updated" json:"updated"`
	Name     string           `bson:"name" json:"name"`
	Versions []ServiceVersion `bson:"versions" json:"versions"`
	Tags     []bson.ObjectId  `bson:"tags" json:"tags"`
}

// ServiceVersion defines a given version of a service
// It contains the version number and all its metadata
// A service version is composed of many containers (with a given image tag)
type ServiceVersion struct {
	Active          bool               `bson:"active" json:"active"` // false means this version should not be used anymore
	Created         time.Time          `bson:"created" json:"created"`
	Updated         time.Time          `bson:"updated" json:"updated"`
	Changelog       string             `bson:"changelog" json:"changelog"`   // Markdown description of what has changed in this version
	Number          string             `bson:"number" json:"number"`         // The version number of the
	PreviousVersion string             `bson:"previous" json:"previous"`     // Previous version of the service
	Containers      []ServiceContainer `bson:"containers" json:"containers"` // 1 or more container defining the service. e.g. a database and the application are two container for 1 service
	// Means automatical update (or notification to update) will not be performed to this version
	// i.e. an old service will only be updatable to the latest version before a breaking change
	HasBreakingChange bool `bson:"hasBreakingChange" json:"hasBreakingChange"`
}

// ServiceContainer is a container of a given image version
// It should host only one process (a database or an independent application)
type ServiceContainer struct {
	Commands Commands `bson:"commands" json:"commands"`
	URLs     URLs     `bson:"urls" json:"urls"`
	Jobs     Jobs     `bson:"jobs" json:"jobs"`
	Image    Image    `bson:"image" json:"image"` // container image
}

// // AddImage adds an Image to the Service
// func (s *Service) AddImage(i *Image) {
// 	s.Images = append(s.Images, *i)
// }

// // AddCommand adds a Command to the Service
// func (s *Service) AddCommand(c *Command) {
// 	s.Commands = append(s.Commands, *c)
// }

// // AddURL adds an URL to the Service
// func (s *Service) AddURL(u *URL) {
// 	s.URLs = append(s.URLs, *u)
// }

// // AddJob adds a Job to the Service
// func (s *Service) AddJob(j *Job) {
// 	s.Jobs = append(s.Jobs, *j)
// }

// // GetLatestImage gets the last created image for given service
// func (s Service) GetLatestImage() (Image, error) {
// 	var last time.Time
// 	var image Image

// 	for _, v := range s.Images {
// 		created := v.Created
// 		if v.Created.After(last) {
// 			last = created
// 			image = v
// 		}
// 	}

// 	if image.Name == "" {
// 		return image, errors.New("Did not find any image")
// 	}

// 	return image, nil
// }

// // GetImage returns the image return from the service
// func (s Service) GetImage(image string) (Image, error) {
// 	for _, v := range s.Images {
// 		if strings.TrimSpace(image) == strings.TrimSpace(v.Name) {
// 			return v, nil
// 		}
// 	}
// 	return Image{}, fmt.Errorf("Did not find image %v in service %v", image, s.Title)
// }

// // IsExistingImage checks that image exists in service
// func (s Service) IsExistingImage(image string) bool {
// 	for _, v := range s.Images {
// 		if strings.TrimSpace(image) == strings.TrimSpace(v.Name) {
// 			return true
// 		}
// 	}

// 	return false
// }

// // GetActiveJobs get active jobs from service.
// func (s Service) GetActiveJobs() (jobs []Job) {
// 	for _, j := range s.Jobs {
// 		if j.Active {
// 			jobs = append(jobs, j)
// 		}
// 	}
// 	return
// }
