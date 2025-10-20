package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"users-api/internal/config"
	"users-api/internal/controllers"
	"users-api/internal/middleware"
	"users-api/internal/repository"
	"users-api/internal/services"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load() // opcional: carga .env si existe

	cfg := config.Load()

	// DB
	db, err := gorm.Open(mysql.Open(cfg.DBDSN), &gorm.Config{})
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
		w.Write([]byte("ok"))
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
