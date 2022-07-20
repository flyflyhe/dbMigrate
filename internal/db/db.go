package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"sync"
)

var dbPoolObj *dbPool
var lock sync.Mutex

type dbPool struct {
	mysql      map[string]*gorm.DB
	sqlserver  map[string]*gorm.DB
	postgresql map[string]*gorm.DB
}

type Operator struct {
	Select func()
	Insert func()
	Delete func()
}

func init() {
	dbPoolObj = &dbPool{
		mysql:      make(map[string]*gorm.DB),
		sqlserver:  make(map[string]*gorm.DB),
		postgresql: make(map[string]*gorm.DB),
	}
}

func GetDb(dsn string, dsnType string) (*gorm.DB, error) {
	var err error

	if dsnType == "mysql" {
		if db, ok := dbPoolObj.mysql[dsn]; ok {
			return db, err
		} else {
			lock.Lock()
			defer lock.Unlock()
			db, err = gorm.Open(mysql.New(mysql.Config{
				DSN:                       dsn,   // data source name
				DefaultStringSize:         10,    // default size for string fields
				DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
				DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
				DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
				SkipInitializeWithVersion: false, // // auto configure based on currently MySQL version
			}))

			dbPoolObj.mysql[dsn] = db
			return db, err
		}
	}

	return nil, nil
}

func CreateTable(db0, db1 *gorm.DB, t0, t1 string) error {
	var result map[string]interface{}
	if err := db0.Raw("show create table " + t0).Take(&result).Error; err != nil {
		return err
	}

	createSql := result["Create Table"].(string)

	createSqlNew := strings.Replace(createSql, t0, t1, 1)

	if err := db1.Exec(createSqlNew).Error; err != nil {
		return err
	}

	return nil
}

func GetId(id interface{}) uint64 {
	var startId uint64
	switch id.(type) {
	case uint64:
		startId = id.(uint64)
	case uint32:
		startId = uint64(id.(uint32))
	case int:
		startId = uint64(id.(int))
	case uint:
		startId = uint64(id.(uint))
	}

	return startId
}
