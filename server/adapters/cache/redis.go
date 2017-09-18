package cache

import (
	"encoding/json"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"gopkg.in/redis.v3"
)

// Cache interface
type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
}

// Redis used as cache
type Redis struct {
	client *redis.Client
}

// NewRedis cache
func NewRedis(client *redis.Client) Cache {
	return &Redis{client}
}

// Get the value from redis if the client is defined
func (r Redis) Get(key string) (interface{}, error) {
	var value interface{}
	if r.client == nil {
		return nil, errors.New("Redis is unavailable")
	}
	str, err := r.client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(str), value)
	return value, err
}

// Set the value into redis if the client exist
func (r Redis) Set(key string, value interface{}, ttl time.Duration) error {
	if r.client == nil {
		return errors.New("Redis is unavailable")
	}
	bytes, err := json.Marshal(value)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"value": value,
		}).Error("Invalid JSON")
		return err
	}
	if err = r.client.Set(key, bytes, ttl).Err(); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"key":   key,
			"value": string(bytes),
		}).Error("Cannot set value in Redis")
		return err
	}
	return nil
}
