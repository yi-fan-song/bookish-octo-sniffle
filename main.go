package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"octo/data"
	"octo/handlers"
	"octo/ticker"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
)

var (
	isProduction = flag.Bool("production", false, "Indicates production environment")
	hostname     = flag.String("hostname", "", "Hostname of the app in production, this will be used to generate a certificate from Let's Encrypt")
	port         = flag.Int("port", 80, "Port to use to serve http")
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	dbService := &data.Service{DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Toronto", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))}

	t := ticker.Ticker{DbService: dbService}
	t.Start()

	if err := dbService.Migrate(); err != nil {
		log.Fatal("failed to migrate models")
	}

	e.GET("/status", handlers.GetStatus)
	e.POST("/status", handlers.GetStatus)

	e.POST("/timer", func(c echo.Context) error { return handlers.PostTimer(c, dbService) })
	e.GET("/timer/:id", func(c echo.Context) error { return handlers.GetTimer(c, dbService) })
	e.PUT("/timer/:id", func(c echo.Context) error { return handlers.PutTimer(c, dbService) })

	if *isProduction {
		if *hostname == "" {
			panic("no hostname specified")
		}

		e.Pre(middleware.HTTPSRedirect())

		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(*hostname)
		go func() { e.Logger.Fatal(e.StartAutoTLS(":443")) }()
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}
