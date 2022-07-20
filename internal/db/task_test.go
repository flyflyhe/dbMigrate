package db

import (
	_ "embed"
	"errors"
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
		created:        true,
	}
	task.SetFuncByConfig(taskConfig)

	log.Println(task.Migrate())
}

func TestTask(t *testing.T) {
	task := CreateTask(dsn1, dsn0, "user_bak", "user", "mysql", "mysql")
	task.SetStart(func() (map[string]interface{}, error) {
		conn, err := task.GetDb(task.Dsn0(), task.DsnType0())
		if err != nil {
			log.Println(err)
			return nil, err
		}

		var result map[string]interface{}

		if err = conn.Table(task.T0()).Limit(1).Take(&result).Error; err != nil {
			log.Println(err)
			return nil, err
		}

		return result, nil
	})

	task.SetNext(func(m map[string]interface{}) (map[string]interface{}, error) {
		if m == nil {
			return nil, errors.New("无数据了")
		}
		conn, err := task.GetDb(task.Dsn0(), task.DsnType0())
		if err != nil {
			log.Println(err)
			return nil, err
		}

		id := m["id"].(uint32)

		var result map[string]interface{}

		if err := conn.Table(task.T0()).Where("id > ?", id).Limit(1).Take(&result).Error; err != nil {
			return nil, err
		}

		return result, nil
	})

	task.SetCreate(func(m map[string]interface{}) error {
		conn, err := task.GetDb(task.Dsn1(), task.DsnType1())
		if err != nil {
			log.Println(err)
			return err
		}
		return conn.Table(task.T1()).Create(&m).Error
	})

	log.Println(task.Migrate())
}
