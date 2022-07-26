package controller

import (
	"encoding/json"
	"github.com/flyflyhe/dbMigrate/internal/db"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestTask_Create(t *testing.T) {
	taskConfig := &db.TaskExternalConfig{
		Name:          "test1",
		Dsn0:          "dsn0",
		Dsn1:          "dsn1",
		T0:            "t0",
		T1:            "t1",
		DsnType0:      "dsnType0",
		DsnType1:      "dsnType1",
		StartFuncType: 2,
		NextFuncType:  2,
		NextKey:       "id",
		EndFuncType:   2,
		EndKey:        "id",
		EndVal:        100,
		DeleteKey:     "",
		Created:       false,
	}
	taskConfigBytes, _ := json.Marshal(taskConfig)
	request, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/tasks", strings.NewReader(string(taskConfigBytes)))
	if err != nil {
		t.Error(err)
	}
	request.Header.Set("Content-Type", "application/json")

	if res, err := http.DefaultClient.Do(request); err != nil {
		t.Error(err)
	} else if res.StatusCode != 200 {
		t.Error(res.StatusCode)
		if result, err := ioutil.ReadAll(res.Body); err != nil {
			t.Error(err)
		} else {
			t.Log(string(result))
		}
	} else {
		if result, err := ioutil.ReadAll(res.Body); err != nil {
			t.Error(err)
		} else {
			t.Log(string(result))
		}
	}
}
