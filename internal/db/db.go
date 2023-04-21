package db

import (
	"gorm.io/gorm"
	"strings"
)

type Wrapper struct {
	*gorm.DB
}

func (w *Wrapper) AllTables(database string) ([]string, error) {
	var result []string
	if err := w.Debug().Raw("select TABLE_NAME from information_schema.tables where table_schema=?", database).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (w *Wrapper) TableSchema(table string) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := w.Debug().Raw("show create table " + table).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (w *Wrapper) CreateTable(ddl string) error {
	if err := w.Debug().Exec(ddl).Error; err != nil {
		return err
	}

	return nil
}

func (w *Wrapper) ChangeDDL(oTable, nTable, ddl string) string {
	ddl = ddl[:strings.LastIndex(ddl, ")")+1]
	return strings.Replace(ddl, oTable, nTable, 1)
}
