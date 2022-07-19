package db

import (
	"gorm.io/gorm"
	"log"
	"sync"
)

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
	wg.Add(2)
	go func() {
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
