package connection

import (
	"fmt"
	"github.com/dbMigrate/v2/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDb(config config.MysqlConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		config.User, config.Password, config.Host, config.Port, config.Database,
	) + "&loc=Asia%2fShanghai"
	fmt.Println("dsn:" + dsn)
	if db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{}); err != nil {
		return nil, err
	} else {
		if sqlDb, err := db.DB(); err != nil {
			return nil, err
		} else {
			sqlDb.SetMaxIdleConns(20)
			sqlDb.SetMaxOpenConns(40)
		}
		return db, nil
	}
}
