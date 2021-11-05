package data

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Service struct {
	DSN string
	db  *gorm.DB
}

func (s *Service) Db() (*gorm.DB, error) {
	if s.db != nil {
		return s.db, nil
	}
	db, err := gorm.Open(postgres.Open(s.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	s.db = db
	return s.db, nil
}

func (s *Service) Migrate() error {
	db, err := s.Db()
	if err != nil {
		return err
	}

	db.AutoMigrate(&Timer{})
	return nil
}
