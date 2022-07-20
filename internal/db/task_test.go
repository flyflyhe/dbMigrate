package db

import (
	_ "embed"
	"log"
	"testing"
)

//go:embed dsn0.txt
var dsn0 string

//go:embed dsn1.txt
var dsn1 string

func TestCreateTask(t *testing.T) {
	task := CreateTask(dsn0, dsn1, "user", "user_bak", "mysql", "mysql")
	taskConfig := &TaskConfig{
		startCondition: map[string][]interface{}{"id >= ?": []interface{}{1}},
		startFuncType:  startFuncTypeCustom,
		nextFuncType:   nextFuncTypeId,
		nextKey:        "id",
		endFuncType:    endFuncTypeId,
		endKey:         "id",
		endVal:         100,
		task:           task,
		created:        true,
	}
	task.SetFuncByConfig(taskConfig)

	log.Println(task.Migrate())
}
