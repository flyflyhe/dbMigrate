package model

import "time"

type Task struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	Name           string    `json:"name" gorm:"unique;not null;"`
	Dsn0           string    `json:"dsn_0"`
	Dsn1           string    `json:"dsn_1"`
	T0             string    `json:"t_0"`
	T1             string    `json:"t_1"`
	DsnType0       string    `json:"dsn_type_0"`
	DsnType1       string    `json:"dsn_type_1"`
	StartCondition string    `json:"start_condition"`
	StartFuncType  int       `json:"start_func_type"`
	NextFuncType   int       `json:"next_func_type"`
	NextKey        string    `json:"next_key"`
	EndFuncType    int       `json:"end_func_type"`
	EndKey         string    `json:"end_key"`
	EndVal         string    `json:"end_val"`
	DeleteKey      string    `json:"delete_key"`
	Created        int       `json:"created"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (task *Task) TableName() string {
	return "task"
}
