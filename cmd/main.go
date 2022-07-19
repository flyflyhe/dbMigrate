package main

import (
	_ "embed"
	"errors"
	"github.com/flyflyhe/dbMigrate/internal/db"
	"log"
)

//go:embed dsn0.txt
var dsn0 string

//go:embed dsn1.txt
var dsn1 string

func main() {
	//err := db.Migrate(dsn0, dsn1, "user", "user_bak", "mysql", "mysql")
	//log.Println(err)

	task := db.CreateTask(dsn1, dsn0, "user_bak", "user", "mysql", "mysql")
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

	log.Println(task.Migrate())
}
