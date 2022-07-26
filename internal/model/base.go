package model

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Base interface {
	TableName() string
	GetDb() (*gorm.DB, error)
}

func Create(m Base) error {
	db, err := m.GetDb()
	if err != nil {
		return err
	}

	log.Debug().Caller().Msg(m.TableName())
	return db.Table(m.TableName()).Create(m).Error
}
