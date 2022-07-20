package db

import (
	"encoding/json"
	"log"
	"testing"
)

func TestGetConfig(t *testing.T) {
	if l, err := GetConfig([]byte(TaskYaml)); err != nil {
		t.Error(err)
	} else {
		m, _ := json.Marshal(l)
		log.Println(string(m))
	}
}
