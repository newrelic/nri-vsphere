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
}

var ReadTimestampCache = readTimestampCache
var WriteTimestampCache = writeTimestampCache

func readTimestampCache(directoryPath string, resourceName string) (time.Time, error) {
	filePath := path.Join(directoryPath, getName(resourceName))

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while opening cache file: %s, ", err.Error())
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while reading cache file: %s, ", err.Error())
	}
	var c Cache
	err = json.Unmarshal(byteValue, &c)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while unmashalling cache file content: %s, ", err.Error())
	}

	timestamp, err := time.Parse(time.RFC1123, c.LastTimestamp)
	return timestamp, err
}

func writeTimestampCache(directoryPath string, resourceName string, lastTimestamp time.Time) error {
	filePath := path.Join(directoryPath, getName(resourceName))

	err := os.Mkdir(directoryPath, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("error whilecreating directory: %s, %s ", err.Error(), directoryPath)
	}

	jsonFile, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	c := Cache{
		LastTimestamp: lastTimestamp.Format(time.RFC1123),
	}
	s, err := json.Marshal(c)
	if err != nil {
		return err
	}
	_, err = jsonFile.WriteAt([]byte(s), 0)
	return err
}

func getName(resource string) string {
	return "vsphere-" + resource + ".json"
}
