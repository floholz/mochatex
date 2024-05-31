package parsing

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func ParseJsonFile(path *string, errLog, infoLog *log.Logger) map[string]interface{} {
	if filepath.Ext(*path) != ".json" {
		errLog.Fatalf("%s must be a valid .json file", *path)
	}
	_, err := os.Stat(*path)
	if err != nil {
		errLog.Fatalf("error while reading info for %s: %v", *path, err)
	}

	var dtls map[string]interface{}
	dFile, err := os.Open(*path)
	if err != nil {
		errLog.Fatalf("error while opening details json file %s: %v", *path, err)
	}
	err = json.NewDecoder(dFile).Decode(&dtls)
	if err != nil {
		errLog.Fatalf("error while decoding json file %s: %v", *path, err)
	}

	return dtls
}
