package tasks

import (
	"github.com/dbMigrate/v2/internal/db"
	"github.com/dbMigrate/v2/internal/scripts"
	"github.com/dbMigrate/v2/pkg/logging"
	"strings"
	"sync"
	"time"
)

type Task struct {
	Source         *db.Wrapper
	SourceDatabase string
	Dst            *db.Wrapper
	DstDatabase    string
	tableList      []string
	ddlMap         map[string]string
	tableColumn    map[string]map[string]db.Columns
}

func (t *Task) InitTables() error {
	if tableList, err := t.Source.AllTables(); err != nil {
		return err
	} else {
		filterResult := make([]string, 0)
		for _, table := range tableList {
			if scripts.Filter(table) == scripts.Allow {
				filterResult = append(filterResult, table)
			}
		}
		t.tableList = filterResult
		return nil
	}
}

func (t *Task) InitColumnDDL() error {
	if t.ddlMap == nil {
		t.ddlMap = make(map[string]string)
	}
	if t.tableColumn == nil {
		t.tableColumn = map[string]map[string]db.Columns{}
	}
	for _, table := range t.tableList {
		t.tableColumn[table] = map[string]db.Columns{}

		//ddl格式转换
		if ddl, err := t.Source.TableSchema(table); err != nil {
			return err
		} else {
			changeDDL := t.Source.ChangeDDL(table, table, ddl)
			t.ddlMap[table] = scripts.Convert(table, changeDDL)
		}

		//列类型保存
		if columns, err := t.Source.TableColumns(t.SourceDatabase, table); err != nil {
			return err
		} else {
			for _, column := range columns {
				t.tableColumn[table][column.ColumnName] = column
			}
		}
	}
	return nil
}

func (t *Task) start() error {
	for table, ddl := range t.ddlMap {
		if err := t.Dst.CreateTable(table, ddl); err != nil {
			return err
		}
	}

	wg := sync.WaitGroup{}
	for _, table := range t.tableList {
		wg.Add(1)
		go func(table string) {
			defer func() {
				wg.Done()
			}()
			dataChan := t.Source.ScanDataByTable(table)

			var result []map[string]interface{}
			for data := range dataChan {
				t.convertColumnValue(table, data)
				result = append(result, data)
				if len(result) == 100 {
					if err := t.Dst.BatchInsert(table, result); err != nil {
						logging.Logger.Sugar().Error(err)
						return
					}
					result = []map[string]interface{}{}
				}
			}

			if len(result) > 0 {
				if err := t.Dst.BatchInsert(table, result); err != nil {
					logging.Logger.Sugar().Error(err)
					return
				}
			}
		}(table)
	}
	wg.Wait()
	return nil
}

func (t *Task) Start() error {
	if err := t.InitTables(); err != nil {
		return err
	}

	if err := t.InitColumnDDL(); err != nil {
		return err
	}

	if err := t.start(); err != nil {
		return err
	}

	return nil
}

func (t *Task) compare() error {
	for _, table := range t.tableList {
		sC, err := t.Source.GetCount(table)
		if err != nil {
			logging.Logger.Sugar().Error(err)
		}

		dC, err := t.Dst.GetCount(table)
		if err != nil {
			logging.Logger.Sugar().Error(err)
		}

		logging.Logger.Sugar().Info("table:", table, "sC:", sC, "dC", dC)
		if sC != dC {
			logging.Logger.Sugar().Error(table, "同步数据少")
		} else {
			logging.Logger.Sugar().Info(table, "同步数据正常")
		}
	}
	return nil
}

func (t *Task) Compare() error {
	if err := t.InitTables(); err != nil {
		return err
	}

	if err := t.InitColumnDDL(); err != nil {
		return err
	}

	if err := t.compare(); err != nil {
		return err
	}

	return nil
}

func (t *Task) convertColumnValue(table string, data map[string]interface{}) {
	column := t.tableColumn[table]
	if table == "log" {
		logging.Logger.Sugar().Info(data)
		logging.Logger.Sugar().Info(column)
	}
	result := data
	for k, v := range data {
		if column[k].IsNullable != "YES" {
			continue
		}
		t := column[k].DataType
		if t == "date" {
			if s, ok := v.(time.Time); ok {
				if strings.Contains(s.String(), "0000-00-00") || strings.Contains(s.String(), "0001-01-01") {
					result[k] = nil
				}
			}
		} else if t == "datetime" {
			if s, ok := v.(time.Time); ok {
				if strings.Contains(s.String(), "0000-00-00 00:00:00") || strings.Contains(s.String(), "0001-01-01 00:00:00") {
					result[k] = nil
				}
			}
		} else if t == "time" {
			if s, ok := v.(time.Time); ok {
				if strings.Contains(s.String(), "00:00:00") || strings.Contains(s.String(), "00:00:00") {
					result[k] = nil
				}
			}
		}
	}
}
