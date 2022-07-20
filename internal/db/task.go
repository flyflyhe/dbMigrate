package db

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"strings"
	"sync"
)

const (
	startFuncTypeCustom  = 1
	startFuncTypeDefault = 2

	nextFuncTypeCustom = 1
	nextFuncTypeId     = 2

	endFuncTypeCustom   = 1
	endFuncTypeId       = 2
	endFuncTypeDatetime = 3
)

type TaskConfig struct {
	startCondition map[string][]interface{}
	startFuncType  int    //default
	nextFuncType   int    //id
	nextKey        string //id
	endFuncType    int    //id  datetime
	endKey         string
	endVal         interface{}
	deleteKey      string
	created        bool
}

func CreateTaskConfigByEConfig(eConfig *TaskExternalConfig) *TaskConfig {
	return &TaskConfig{
		startCondition: eConfig.StartCondition,
		startFuncType:  eConfig.StartFuncType,
		nextFuncType:   eConfig.NextFuncType,
		nextKey:        eConfig.NextKey,
		endFuncType:    eConfig.EndFuncType,
		endKey:         eConfig.EndKey,
		endVal:         eConfig.EndVal,
		deleteKey:      eConfig.DeleteKey,
		created:        eConfig.Created,
	}
}

func (task *Task) SetFuncByConfig(config *TaskConfig) {
	task.SetStart(func() (map[string]interface{}, error) {
		conn, err := task.GetDb(task.Dsn0(), task.DsnType0())
		if err != nil {
			log.Println(err)
			return nil, err
		}

		var result map[string]interface{}

		if config.startFuncType == startFuncTypeDefault {
			if err = conn.Table(task.T0()).Limit(1).Debug().Take(&result).Error; err != nil {
				log.Println(err)
				return nil, err
			}
		} else if config.startFuncType == startFuncTypeCustom {
			var where string
			var val []interface{}
			for where, val = range config.startCondition {
				break
			}
			if err = conn.Table(task.T0()).Where(where, val...).Limit(1).Debug().Take(&result).Error; err != nil {
				log.Println(err)
				return nil, err
			}
		}

		log.Println(result)
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

		var result map[string]interface{}

		if config.nextFuncType == nextFuncTypeId {
			if err := conn.Table(task.T0()).Where(config.nextKey+" > ?", m[config.nextKey]).Limit(1).Take(&result).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("暂不支持")
		}

		log.Println("next", result)
		return result, nil
	})

	if config.endFuncType != 0 {
		task.SetEnd(func(m map[string]interface{}) bool {
			if config.endFuncType == endFuncTypeId {
				endVal := GetId(config.endVal)
				id := GetId(m[config.endKey])
				return id >= endVal
			} else if config.startFuncType == endFuncTypeDatetime {
				endVal := config.endVal.(string)
				return strings.Compare(endVal, m[config.endKey].(string)) >= 0
			}
			return false
		})
	}

	if config.created {
		task.SetCreate(func(m map[string]interface{}) error {
			log.Println("create", m)
			conn, err := task.GetDb(task.Dsn1(), task.DsnType1())
			if err != nil {
				log.Println(err)
				return err
			}
			return conn.Table(task.T1()).Create(&m).Error
		})
	}

	if config.deleteKey != "" {
		task.SetDelete(func(m map[string]interface{}) error {
			if m == nil {
				return errors.New("无数据")
			}
			conn, err := task.GetDb(task.Dsn0(), task.DsnType0())
			if err != nil {
				log.Println(err)
				return err
			}

			return conn.Table(task.t0).Where(config.deleteKey+" = ?", m[config.deleteKey]).Delete(m).Error
		})
	}
}

type Task struct {
	dsn0     string
	dsn1     string
	t0       string
	t1       string
	dsnType0 string
	dsnType1 string
	start    func() (map[string]interface{}, error)                       //开始方法
	end      func(map[string]interface{}) bool                            //结束判断
	next     func(map[string]interface{}) (map[string]interface{}, error) //迭代方法
	delete   func(map[string]interface{}) error                           //delete 方法 未设置则不删除
	create   func(map[string]interface{}) error                           //创建方法 不设置则不创建 控制变量可以实现先迁移再删除
}

func (task *Task) Dsn0() string {
	return task.dsn0
}

func (task *Task) SetDsn0(dsn0 string) {
	task.dsn0 = dsn0
}

func (task *Task) Dsn1() string {
	return task.dsn1
}

func (task *Task) SetDsn1(dsn1 string) {
	task.dsn1 = dsn1
}

func (task *Task) T0() string {
	return task.t0
}

func (task *Task) SetT0(t0 string) {
	task.t0 = t0
}

func (task *Task) T1() string {
	return task.t1
}

func (task *Task) SetT1(t1 string) {
	task.t1 = t1
}

func (task *Task) DsnType0() string {
	return task.dsnType0
}

func (task *Task) SetDsnType0(dsnType0 string) {
	task.dsnType0 = dsnType0
}

func (task *Task) DsnType1() string {
	return task.dsnType1
}

func (task *Task) SetDsnType1(dsnType1 string) {
	task.dsnType1 = dsnType1
}

func CreateTask(dsn0, dsn1, t0, t1, dsnType0, dsnTyp1 string) *Task {
	return &Task{
		dsn0:     dsn0,
		dsn1:     dsn1,
		t0:       t0,
		t1:       t1,
		dsnType0: dsnType0,
		dsnType1: dsnTyp1,
	}
}

func (task *Task) SetStart(start func() (map[string]interface{}, error)) {
	task.start = start
}

func (task *Task) SetEnd(end func(map[string]interface{}) bool) {
	task.end = end
}

func (task *Task) SetNext(next func(map[string]interface{}) (map[string]interface{}, error)) {
	task.next = next
}

func (task *Task) SetDelete(delete func(map[string]interface{}) error) {
	task.delete = delete
}

func (task *Task) SetCreate(create func(map[string]interface{}) error) {
	task.create = create
}

func (task *Task) GetDb(dsn, dsnType string) (*gorm.DB, error) {
	return GetDb(dsn, dsnType)
}

func (task *Task) checkTable() error {
	db0, err := task.GetDb(task.dsn0, task.dsnType0)
	if err != nil {
		return err
	}
	db1, err := task.GetDb(task.dsn1, task.dsnType1)
	if err != nil {
		return err
	}
	if !db1.Migrator().HasTable(task.t1) {
		if err := CreateTable(db0, db1, task.t0, task.t1); err != nil {
			return err
		}
	}
	return nil
}

func (task *Task) Migrate() error {
	if err := task.checkTable(); err != nil {
		return err
	}

	resultChan := make(chan map[string]interface{}, 100)
	deleteChan := make(chan map[string]interface{}, 100)
	defer func() {
		close(resultChan)
		close(deleteChan)
	}()

	stopFunc := func() {
		resultChan <- nil
		deleteChan <- nil
	}

	sendResult := func(result map[string]interface{}) {
		resultChan <- result
		deleteChan <- result
	}
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		var err error
		var result map[string]interface{}
		result, err = task.start() //start 返回数据 作为next的条件
		if err != nil {
			stopFunc()
			return
		}

		sendResult(result)

		for {
			if result, err = task.next(result); err != nil {
				log.Println(err) //无数据时会输出错误
				stopFunc()
				break
			} else {
				if task.end != nil && task.end(result) {
					log.Println("end", result)
					stopFunc()
					break
				}
				sendResult(result)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case result := <-resultChan:
				if result == nil { //收到结束信息
					return
				}
				if task.create != nil {
					if err := task.create(result); err != nil {
						log.Println(err)
					}
				}
			default:
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case result := <-deleteChan:
				if result == nil { //收到结束信息
					return
				}
				if task.delete != nil {
					if err := task.delete(result); err != nil {
						log.Println(err)
					}
				}
			default:
			}
		}
	}()

	wg.Wait()

	return nil
}
