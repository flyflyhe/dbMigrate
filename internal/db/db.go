package db

import (
	"errors"
	"github.com/dbMigrate/v2/pkg/logging"
	"gorm.io/gorm"
	"strings"
)

type Wrapper struct {
	*gorm.DB
}

func (w *Wrapper) AllTables() ([]string, error) {
	return w.Debug().Migrator().GetTables()
}

func (w *Wrapper) TableSchema(table string) (string, error) {
	var result map[string]interface{}
	if err := w.Debug().Raw("show create table " + table).Scan(&result).Error; err != nil {
		return "", err
	}

	if ddl, ok := result["Create Table"]; ok {
		return ddl.(string), nil
	} else {
		return "", errors.New(table + "不存在")
	}
}

func (w *Wrapper) CreateTable(table, ddl string) error {
	if w.Debug().Migrator().HasTable(table) {
		return nil
	}
	if err := w.Debug().Exec(ddl).Error; err != nil {
		return err
	}

	return nil
}

func (w *Wrapper) ChangeDDL(oTable, nTable, ddl string) string {
	ddl = ddl[:strings.LastIndex(ddl, ")")+1]
	return strings.Replace(ddl, oTable, nTable, 1)
}

func (w *Wrapper) ScanDataByTable(table string) chan map[string]interface{} {
	dataChan := make(chan map[string]interface{}, 1)
	go func() {
		defer func() {
			close(dataChan)
		}()

		rows, err := w.Debug().Table(table).Rows()
		if err != nil {
			return
		}

		for rows.Next() {
			var result map[string]interface{}
			if err = w.ScanRows(rows, &result); err != nil {
				logging.Logger.Sugar().Error(err)
				return
			} else {
				dataChan <- result
			}
		}
	}()

	return dataChan
}

func (w *Wrapper) BatchInsert(table string, data []map[string]interface{}) error {
	return w.Debug().Table(table).CreateInBatches(data, 100).Error
}
