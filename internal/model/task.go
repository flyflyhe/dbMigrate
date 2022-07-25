package model

import "time"

type Task struct {
	ID             int64  `gorm:"primaryKey"`
	Name           string `gorm:"unique;not null;"`
	Dsn0           string
	Dsn1           string
	T0             string
	T1             string
	DsnType0       string
	DsnType1       string
	StartCondition string
	StartFuncType  int
	NextFuncType   int
	NextKey        string
	EndFuncType    int
	EndKey         string
	EndVal         string
	DeleteKey      string
	Created        int
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (task *Task) TableName() string {
	return "task"
}
