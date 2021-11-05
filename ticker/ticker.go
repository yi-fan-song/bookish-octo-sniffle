package ticker

import (
	"fmt"
	"net/http"
	"octo/data"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

const failedRequestAllowedCount = 3

type Ticker struct {
	DbService *data.Service
}

func (t *Ticker) Start() {
	go func() {
		for {
			time.Sleep(time.Second)
			timers, err := t.DbService.GetExpiredTimers()
			if err != nil {
				log.Errorf("Error getting expired timers: %s", err)
				continue
			}

			for _, timer := range timers {
				go func(timer data.Timer) {
					err := t.callWebhook(timer)
					if err != nil {
						log.Errorf("Error calling webhook url for timer id=%d: %s", timer.ID, err)
						if timer.FailedRequests >= failedRequestAllowedCount {
							// Don't waste ressources trying to call the failing webhook again
							t.DbService.DeleteTimer(timer.ID)
						}
						t.DbService.IncrementFailedRequests(timer.ID)
					}
					t.DbService.DeleteTimer(timer.ID)
				}(timer)
			}
		}
	}()
}

func (t *Ticker) callWebhook(timer data.Timer) error {
	req, err := http.NewRequest(http.MethodPost, timer.WebhookUrl, strings.NewReader("{\"message\": \"Time's up!\"}"))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}
	return nil
}
