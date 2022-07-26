package taskModel

import (
	"encoding/json"
	"fmt"
	"github.com/flyflyhe/dbMigrate/internal/db"
	"github.com/flyflyhe/dbMigrate/internal/store"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type Task struct {
	ID             int64                  `json:"id" gorm:"primaryKey"`
	Name           string                 `json:"name" gorm:"unique;not null;"`
	Dsn0           string                 `json:"dsn_0"`
	Dsn1           string                 `json:"dsn_1"`
	T0             string                 `json:"t_0"`
	T1             string                 `json:"t_1"`
	DsnType0       string                 `json:"dsn_type_0"`
	DsnType1       string                 `json:"dsn_type_1"`
	Condition      string                 `json:"condition"`
	ExternalConfig *db.TaskExternalConfig `gorm:"-:all"`
	CreatedAt      time.Time              `gorm:"autoCreateTime"`
	UpdatedAt      time.Time              `gorm:"autoUpdateTime"`
}

func init() {
	task := &Task{}
	fmt.Println("task init")
	if conn, err := task.GetDb(); err != nil {
		fmt.Println(err.Error())
	} else {
		if !conn.Migrator().HasTable(task.TableName()) {
			if err := conn.Migrator().AutoMigrate(task); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func CreateTaskByConfig(config *db.TaskExternalConfig) *Task {
	conditionBytes, _ := json.Marshal(config)
	return &Task{
		Name:           config.Name,
		Dsn0:           config.Dsn0,
		Dsn1:           config.Dsn1,
		T0:             config.T0,
		T1:             config.T1,
		DsnType0:       config.DsnType0,
		DsnType1:       config.DsnType1,
		Condition:      string(conditionBytes),
		ExternalConfig: config,
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
	}
}

func (task *Task) TableName() string {
	return "task"
}

func (task *Task) GetDb() (*gorm.DB, error) {
	return store.GetDefaultDb()
}

func (task *Task) PrimaryKey() int64 {
	return task.ID
}

func (task *Task) GetConfig() (*db.TaskExternalConfig, error) {
	config := &db.TaskExternalConfig{}

	if err := json.Unmarshal([]byte(task.Condition), config); err != nil {
		return nil, err
	}

	return config, nil
}

func (task *Task) AfterFind(tx *gorm.DB) (err error) {
	if task.Condition != "" {
		var config = &db.TaskExternalConfig{}
		if err = json.Unmarshal([]byte(task.Condition), config); err == nil {
			task.ExternalConfig = config
		} else {
			log.Error().Caller().Err(err).Send()
		}
	}
	return
}
