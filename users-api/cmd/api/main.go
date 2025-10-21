package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"users-api/internal/config"
	"users-api/internal/controllers"
	"users-api/internal/middleware"
	"users-api/internal/repository"
	"users-api/internal/services"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connectWithRetry(dsn string, attempts int, wait time.Duration) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	for i := 1; i <= attempts; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		log.Printf("db connect failed (try %d/%d): %v", i, attempts, err)
		time.Sleep(wait)
	}
	return nil, err
}

func main() {
	_ = godotenv.Load() // en Docker las env vienen del compose
	cfg := config.Load()

	// DB (retry)
	db, err := connectWithRetry(cfg.DBDSN, 30, 1*time.Second)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	if err := repository.AutoMigrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// wiring
	userRepo := repository.NewUserRepo(db)
	authSvc := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTTTL)
	authCtl := controllers.NewAuthController(authSvc)
	usersCtl := controllers.NewUsersController(userRepo)
	authMW := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Gin router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.RecoverJSON())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // poné tu front acá
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// públicos
	r.GET("/health", func(c *gin.Context) { c.String(200, "ok") })
	r.POST("/auth/login", authCtl.Login)
	r.POST("/users", authCtl.RegisterNormal)

	// protegidos
	secured := r.Group("/", authMW.RequireAuth())
	{
		secured.GET("/me", usersCtl.Me)
		secured.GET("/users/:id", usersCtl.GetByID)
		// ejemplo admin-only:
		// admin := secured.Group("/admin", authMW.RequireAdmin())
		// admin.POST("/users", usersCtl.CreateAdmin)
	}

	addr := ":" + cfg.HTTPPort
	log.Printf("users-api (gin) up on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
