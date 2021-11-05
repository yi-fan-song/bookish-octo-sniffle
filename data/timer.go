package data

import (
	"errors"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
)

type Timer struct {
	gorm.Model

	Duration       uint64
	TimeRemaining  uint64
	ExpiresOn      time.Time
	IsPaused       bool
	WebhookUrl     string
	FailedRequests uint64
}

func (s *Service) AddTimer(duration uint64, webhookUrl string) (*Timer, error) {
	timer := Timer{
		Duration:      duration,
		TimeRemaining: duration,
		ExpiresOn:     time.Now().Add(time.Duration(duration) * time.Second),
		IsPaused:      false,
		WebhookUrl:    webhookUrl,
	}

	db, err := s.Db()
	if err != nil {
		return nil, err
	}

	result := db.Create(&timer)
	if result.Error != nil {
		return nil, result.Error
	}

	return &timer, nil
}

var ErrTimerNotFound = fmt.Errorf("timer not found")

func (s *Service) GetTimer(id uint) (*Timer, error) {
	db, err := s.Db()
	if err != nil {
		return nil, err
	}

	timer := Timer{}
	result := db.First(&timer, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrTimerNotFound
		}
		return nil, result.Error
	}

	if !timer.IsPaused {
		timer.TimeRemaining = uint64(math.Round(time.Until(timer.ExpiresOn).Seconds()))

		result = db.Save(timer)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	return &timer, nil
}

func (s *Service) SetTimerPauseStatus(id uint, isPaused bool) (*Timer, error) {
	db, err := s.Db()
	if err != nil {
		return nil, err
	}

	timer := Timer{}
	result := db.First(&timer, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrTimerNotFound
		}
		return nil, result.Error
	}

	if timer.IsPaused != isPaused {
		if isPaused {
			pauseTimer(&timer)
		} else {
			unpauseTimer(&timer)
		}
	}

	result = db.Save(timer)
	if result.Error != nil {
		return nil, result.Error
	}

	return &timer, nil
}

func pauseTimer(timer *Timer) {
	timer.TimeRemaining = uint64(math.Round(time.Until(timer.ExpiresOn).Seconds()))
	timer.IsPaused = true
}

func unpauseTimer(timer *Timer) {
	timer.ExpiresOn = time.Now().Add(time.Duration(timer.TimeRemaining) * time.Second)
	timer.IsPaused = false
}

func (s *Service) GetExpiredTimers() ([]Timer, error) {
	db, err := s.Db()
	if err != nil {
		return nil, err
	}

	timers := []Timer{}
	result := db.Where("expires_on <= ?", time.Now()).Find(&timers)
	if result.Error != nil {
		return nil, result.Error
	}

	return timers, nil
}

func (s *Service) GetTimers(isPaused bool) ([]Timer, error) {
	db, err := s.Db()
	if err != nil {
		return nil, err
	}

	timers := []Timer{}
	result := db.Where(map[string]interface{}{"isPaused": isPaused}).Find(&timers)
	if result.Error != nil {
		return nil, result.Error
	}

	return timers, nil
}

func (s *Service) IncrementFailedRequests(id uint) error {
	db, err := s.Db()
	if err != nil {
		return err
	}

	timer := Timer{}
	result := db.First(&timer, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrTimerNotFound
		}
		return result.Error
	}

	timer.FailedRequests++

	if err := db.Save(&timer).Error; err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteTimer(id uint) error {
	db, err := s.Db()
	if err != nil {
		return err
	}

	timer := Timer{}
	result := db.First(&timer, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrTimerNotFound
		}
		return result.Error
	}

	db.Delete(&timer)
	return nil
}
