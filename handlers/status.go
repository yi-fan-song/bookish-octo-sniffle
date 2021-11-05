package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetStatus(c echo.Context) error {
	fmt.Printf("pinged at %s\n", time.Now().String())
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "pong"})
}
