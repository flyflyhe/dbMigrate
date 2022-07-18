package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
)

var dbPoolObj *dbPool

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

func Migrate(dsn0, dsn1 string, t0, t1 string, dsnType0, dsnType1 string) error {
	db0, err := GetDb(dsn0, dsnType0)
	if err != nil {
		return err
	}
	db1, err := GetDb(dsn1, dsnType1)
	if err != nil {
		return err
	}
	if !db1.Migrator().HasTable(t1) {
		if err := CreateTable(db0, db1, t0, t1); err != nil {
			return err
		}
	}

	migrate(db0, db1, t0, t1)
	return nil
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

func migrate(db0, db1 *gorm.DB, t0, t1 string) {
	resultChan := make(chan map[string]interface{}, 100)
	deleteChan := make(chan uint64, 100)
	stopChan := make(chan struct{})
	defer func() {
		close(resultChan)
		close(deleteChan)
		close(stopChan)
	}()
	stopIndex := 2
	go func() {
		var result map[string]interface{}
		var startId uint64
		if err := db0.Table(t0).Limit(1).Take(&result).Error; err != nil {
			log.Println(err)
			return
		} else {
			if id, ok := result["id"]; !ok {
				log.Println("不支持 无id自增表")
			} else {
				startId = GetId(id)
			}
		}
		for {
			if err := db0.Table(t0).Where("id >= ?", startId).Limit(1).Take(&result).Error; err != nil {
				log.Println(err) //无数据时会输出错误
				log.Println(result)
				resultChan <- nil
				deleteChan <- 0
				break
			} else {
				resultChan <- result
				deleteChan <- startId
				startId = GetId(result["id"]) + 1
			}
		}
	}()

	go func() {
		for {
			select {
			case result := <-resultChan:
				if result == nil { //收到结束信息
					stopChan <- struct{}{}
					return
				}
				if err := db1.Table(t1).Create(&result).Error; err != nil {
					log.Println(err)
				}
			default:
			}
		}
	}()

	go func() {
		for {
			select {
			case id := <-deleteChan:
				if id == 0 { //收到结束信息
					stopChan <- struct{}{}
					return
				}
				if err := db0.Exec("delete from `"+t0+"` where id = ?", id).Error; err != nil {
					log.Println(err)
				}
			default:
			}
		}
	}()

	for {
		select {
		case <-stopChan:
			if stopIndex == 1 {
				return
			} else {
				stopIndex--
			}
		default:
		}
	}
}

func GetId(id interface{}) uint64 {
	var startId uint64
	switch id.(type) {
	case uint64:
		startId = id.(uint64)
	case uint32:
		startId = uint64(id.(uint32))
	}

	return startId
}
