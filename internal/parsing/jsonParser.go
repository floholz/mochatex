package parsing

import (
	"encoding/json"
	"log"
	"maps"
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

func FlattenJson(dtls map[string]interface{}) map[string]string {
	return flattenJsonR("", dtls, make(map[string]string))
}

func flattenJsonR(p string, m map[string]interface{}, res map[string]string) map[string]string {
	for key, value := range m {
		fullP := p + "." + key
		switch value.(type) {
		case string:
			res[fullP] = value.(string)
		default:
			inner := flattenJsonR(fullP, value.(map[string]interface{}), res)
			maps.Copy(res, inner)
		}
	}
	return res
}
