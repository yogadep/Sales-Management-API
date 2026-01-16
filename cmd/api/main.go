package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"sales-management-api/internal/config"
	"sales-management-api/internal/db"

	authH "sales-management-api/internal/handlers/auth"
	"sales-management-api/internal/middlewares"
	"sales-management-api/internal/repositories"
	"sales-management-api/internal/services"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	// _ = database // nanti dipakai step berikutnya (repo/service/handler)

	userRepo := repositories.NewUserRepo(database)
	authSvc := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := authH.New(authSvc)

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
			"env":    cfg.AppEnv,
		})
	})

	api := e.Group("/api")

	api.POST("/login", authHandler.Login)

	protected := api.Group("")
	protected.Use(middlewares.JWTAuth(cfg.JWTSecret))

	protected.POST("/register", authHandler.Register, middlewares.RequireRole("admin"))

	log.Printf("server running on :%s", cfg.AppPort)
	e.Logger.Fatal(e.Start(":" + cfg.AppPort))
}
