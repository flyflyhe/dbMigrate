package store

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
)

var defaultManger *DbManger

type DbManger struct {
	path string
	db   *gorm.DB
	lock sync.Mutex
}

func init() {
	defaultManger = &DbManger{path: "./migrate.db"}
}

func GetDefaultDb() (*gorm.DB, error) {
	var err error
	if defaultManger.db == nil {
		defaultManger.lock.Lock()
		defer defaultManger.lock.Unlock()
		defaultManger.db, err = gorm.Open(sqlite.Open(defaultManger.path), &gorm.Config{})
	}
	return defaultManger.db, err
}
