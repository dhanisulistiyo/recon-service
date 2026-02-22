package main

import (
	"log"

	"reconciliation-service/internal/di"
	"reconciliation-service/internal/shared/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler := di.Initialize()
	handler.Register(e)

	log.Println("Starting", cfg.ApplicationName, "on port", cfg.HTTPPort)
	e.Logger.Fatal(e.Start(":" + cfg.HTTPPort))
}
