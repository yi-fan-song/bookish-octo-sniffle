package handlers

import (
	"errors"
	"net/http"
	"octo/data"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type timerPostRequest struct {
	Duration   uint64 `json:"duration" validate:"required"`
	WebhookUrl string `json:"webhookUrl" validate:"required"`
}

type timerPutRequest struct {
	IsPaused *bool `json:"isPaused"`
}

type timerModel struct {
	Id            string `json:"id"`
	Duration      uint64 `json:"duration"`
	TimeRemaining uint64 `json:"timeRemaining"`
	ExpiresOn     string `json:"expiresOn"`
	IsPaused      bool   `json:"isPaused"`
	WebhookUrl    string `json:"webhookUrl"`
}

// POST "/timer"
func PostTimer(c echo.Context, dbService *data.Service) error {
	req := timerPostRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	timer, err := dbService.AddTimer(req.Duration, req.WebhookUrl)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, timerModel{
		Id:            strconv.FormatUint(uint64(timer.ID), 10),
		Duration:      timer.Duration,
		TimeRemaining: timer.TimeRemaining,
		ExpiresOn:     timer.ExpiresOn.Format(time.RFC822),
		IsPaused:      false,
		WebhookUrl:    timer.WebhookUrl,
	})
}

// GET "/timer/:id"
func GetTimer(c echo.Context, dbService *data.Service) error {
	id := c.Param("id")

	timerId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid timer id", "moreInfo": err.Error()})
	}

	timer, err := dbService.GetTimer(uint(timerId))
	if errors.Is(err, data.ErrTimerNotFound) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "timer not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// return c.JSON(http.StatusOK, timer)
	return c.JSON(http.StatusOK, timerModel{
		Id:            strconv.FormatUint(uint64(timer.ID), 10),
		Duration:      timer.Duration,
		TimeRemaining: timer.TimeRemaining,
		ExpiresOn:     timer.ExpiresOn.Format(time.RFC822),
		IsPaused:      timer.IsPaused,
		WebhookUrl:    timer.WebhookUrl,
	})
}

// PUT "/timer/:id"
func PutTimer(c echo.Context, dbService *data.Service) error {
	id := c.Param("id")

	timerId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid timer id", "moreInfo": err.Error()})
	}

	req := timerPutRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if req.IsPaused == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing isPaused field"})
	}

	timer, err := dbService.SetTimerPauseStatus(uint(timerId), *req.IsPaused)
	if errors.Is(err, data.ErrTimerNotFound) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "timer not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, timerModel{
		Id:            strconv.FormatUint(uint64(timer.ID), 10),
		Duration:      timer.Duration,
		TimeRemaining: timer.TimeRemaining,
		ExpiresOn:     timer.ExpiresOn.Format(time.RFC822),
		IsPaused:      timer.IsPaused,
		WebhookUrl:    timer.WebhookUrl,
	})
}
