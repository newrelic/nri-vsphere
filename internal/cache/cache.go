package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type Cache struct {
	LastTimestamp string `json:"lastTimestamp"`
	resourceName  string
	cachePath     string
}

type CacheInterface interface {
	ReadTimestampCache() (time.Time, error)
	WriteTimestampCache(lastTimestamp time.Time) error
}

func NewCache(resourceName string, cachePath string) *Cache {
	return &Cache{
		resourceName: resourceName,
		cachePath:    cachePath,
	}
}

func (c *Cache) ReadTimestampCache() (time.Time, error) {
	filePath := path.Join(c.cachePath, getName(c.resourceName))

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while opening cache file: %s, ", err.Error())
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while reading cache file: %s, ", err.Error())
	}
	var tmpC Cache
	err = json.Unmarshal(byteValue, &tmpC)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while unmashalling cache file content: %s, ", err.Error())
	}
	c.LastTimestamp = tmpC.LastTimestamp

	timestamp, err := time.Parse(time.RFC1123, c.LastTimestamp)
	return timestamp, err
}

func (c *Cache) WriteTimestampCache(lastTimestamp time.Time) error {
	filePath := path.Join(c.cachePath, getName(c.resourceName))

	err := os.Mkdir(c.cachePath, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("error whilecreating directory: %s, %s ", err.Error(), c.cachePath)
	}

	jsonFile, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	c.LastTimestamp = lastTimestamp.Format(time.RFC1123)
	s, err := json.Marshal(c)
	if err != nil {
		return err
	}
	_, err = jsonFile.WriteAt(s, 0)
	return err
}

func getName(resource string) string {
	return "vsphere-" + resource + ".json"
}
