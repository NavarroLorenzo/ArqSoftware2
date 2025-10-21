package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"users-api/internal/config"
	"users-api/internal/controllers"
	"users-api/internal/middleware"
	"users-api/internal/repository"
	"users-api/internal/services"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// intenta conectar varias veces a MySQL (útil en Docker cuando MySQL aún no levantó)
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
	// carga .env en dev; en Docker las env vienen del compose
	_ = godotenv.Load()

	cfg := config.Load()

	// DB con reintentos (30 intentos x 1s ≈ 30s)
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

	// router + middlewares
	mux := http.NewServeMux()
	withRecover := middleware.RecoverAndJSONErrors
	authMW := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// públicos
	mux.Handle("GET /health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	mux.Handle("POST /auth/login", withRecover(http.HandlerFunc(authCtl.Login)))
	mux.Handle("POST /users", withRecover(http.HandlerFunc(authCtl.RegisterNormal)))

	// protegidos
	mux.Handle("GET /me", withRecover(authMW.RequireAuth(http.HandlerFunc(usersCtl.Me))))
	mux.Handle("GET /users/{id}", withRecover(authMW.RequireAuth(http.HandlerFunc(usersCtl.GetByID))))

	addr := ":" + cfg.HTTPPort
	log.Printf("users-api up on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
