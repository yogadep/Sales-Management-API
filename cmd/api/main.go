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

	productH "sales-management-api/internal/handlers/product"
	reportH "sales-management-api/internal/handlers/report"
	salesH "sales-management-api/internal/handlers/sales"
	userH "sales-management-api/internal/handlers/user"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	// auth repo
	userRepo := repositories.NewUserRepo(database)
	authSvc := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := authH.New(authSvc)

	// user repo
	userSvc := services.NewUserService(userRepo)
	userHandler := userH.New(userSvc)

	// product repo
	productRepo := repositories.NewProductRepo(database)
	productSvc := services.NewProductService(productRepo)
	productHandler := productH.New(productSvc)

	// sale repo
	saleRepo := repositories.NewSaleRepo(database)
	saleSvc := services.NewSaleService(saleRepo)
	saleHandler := salesH.New(saleSvc)

	// report repo
	reportRepo := repositories.NewReportRepo(database)
	reportSvc := services.NewReportService(reportRepo)
	reportHandler := reportH.New(reportSvc)

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// CORS (allow frontend access)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
		},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
		},
		AllowHeaders: []string{
			echo.HeaderAuthorization,
			echo.HeaderContentType,
		},
	}))

	// Security headers (XSS, clickjacking, sniffing)
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "DENY",
		ReferrerPolicy:     "same-origin",
	}))

	// Rate limit (20 req/sec per IP)
	e.Use(middleware.RateLimiter(
		middleware.NewRateLimiterMemoryStore(20),
	))

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

	// auth routes
	protected.POST("/register", authHandler.Register, middlewares.RequireRole("admin"))

	// user routes
	protected.GET("/users", userHandler.List, middlewares.RequireRole("admin"))
	protected.GET("/users/:id", userHandler.Detail, middlewares.RequireRole("admin"))

	// product routes
	protected.GET("/products", productHandler.List) // admin + kasir

	protected.POST("/products", productHandler.Create,
		middlewares.RequireRole("admin"),
	)
	protected.PUT("/products/:id", productHandler.Update,
		middlewares.RequireRole("admin"),
	)
	protected.DELETE("/products/:id", productHandler.Delete,
		middlewares.RequireRole("admin"),
	)

	// sale routes
	protected.POST("/sales", saleHandler.Create, middlewares.RequireRole("admin", "kasir"))
	protected.GET("/sales", saleHandler.List, middlewares.RequireRole("admin", "kasir"))
	protected.GET("/sales/:id", saleHandler.Detail, middlewares.RequireRole("admin", "kasir"))
	protected.POST("/sales", saleHandler.Create, middlewares.RequireRole("admin", "kasir"))

	// report routes
	protected.GET("/reports/sales", reportHandler.SalesJSON,
		middlewares.RequireRole("admin", "kasir"))

	protected.GET("/reports/sales.pdf", reportHandler.SalesPDF,
		middlewares.RequireRole("admin", "kasir"))

	protected.GET("/reports/sales.xlsx", reportHandler.SalesExcel,
		middlewares.RequireRole("admin", "kasir"))

	log.Printf("server running on :%s", cfg.AppPort)
	e.Logger.Fatal(e.Start(":" + cfg.AppPort))
}
