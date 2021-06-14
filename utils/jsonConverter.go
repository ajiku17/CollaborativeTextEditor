package utils

import (
	"encoding/json"
	"fmt"
)

func ToJson(obj interface{}) []byte {
	jsonString, err := json.Marshal(obj)
	if err != nil {
		fmt.Errorf("Couldn't parse to Json")
	}
	return jsonString
}

func FromJson(jsonString []byte, obj interface{}) interface{} {
	json.Unmarshal(jsonString, &obj)
	return obj
}