package model

import (
	"gorm.io/gorm"
)

type Base interface {
	TableName() string
	GetDb() (*gorm.DB, error)
	PrimaryKey() int64
}

func Create(m Base) error {
	db, err := m.GetDb()
	if err != nil {
		return err
	}

	return db.Table(m.TableName()).Create(m).Error
}

func Delete(m Base) error {
	db, err := m.GetDb()
	if err != nil {
		return err
	}

	return db.Table(m.TableName()).Delete(m, m.PrimaryKey()).Error
}

func Update(m Base) error {
	db, err := m.GetDb()
	if err != nil {
		return err
	}

	return db.Table(m.TableName()).Updates(m).Error
}

func FindOne(m Base, id int64) error {
	db, err := m.GetDb()
	if err != nil {
		return err
	}

	return db.Table(m.TableName()).Where("id = ?", id).Find(m).Error
}
