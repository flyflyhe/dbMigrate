package tasks

import (
	"github.com/dbMigrate/v2/internal/db"
	"github.com/dbMigrate/v2/pkg/logging"
	"sync"
)

type Task struct {
	Source         *db.Wrapper
	SourceDatabase string
	Dst            *db.Wrapper
	DstDatabase    string
	tableList      []string
	ddlMap         map[string]string
}

func (t *Task) InitTables() error {
	if tableList, err := t.Source.AllTables(); err != nil {
		return err
	} else {
		t.tableList = tableList
		return nil
	}
}

func (t *Task) InitDDL() error {
	if t.ddlMap == nil {
		t.ddlMap = make(map[string]string)
	}
	for _, table := range t.tableList {
		if ddl, err := t.Source.TableSchema(table); err != nil {
			return err
		} else {
			t.ddlMap[table] = t.Source.ChangeDDL(table, table, ddl)
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

	if err := t.InitDDL(); err != nil {
		return err
	}

	if err := t.start(); err != nil {
		return err
	}

	return nil
}
