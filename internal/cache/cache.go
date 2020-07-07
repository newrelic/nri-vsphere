package cache

import (
	"fmt"
	"github.com/newrelic/infra-integrations-sdk/persist"
	"time"
)

type Cache struct {
	resourceName string
	store        persist.Storer
}

type CacheInterface interface {
	ReadTimestampCache() (time.Time, error)
	WriteTimestampCache(lastTimestamp time.Time) error
}

func NewCache(resourceName string, store persist.Storer) *Cache {
	return &Cache{
		resourceName: resourceName,
		store:        store,
	}
}

func (c *Cache) ReadTimestampCache() (time.Time, error) {
	var ts int64
	_, err := c.store.Get(c.resourceName, &ts)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while reading cache file: %s, ", err.Error())
	}
	return time.Unix(0, ts), err
}

func (c *Cache) WriteTimestampCache(lastTimestamp time.Time) error {
	c.store.Set(c.resourceName, lastTimestamp.UnixNano())
	err := c.store.Save()
	return err
}
