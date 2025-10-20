package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPPort  string
	DBDSN     string
	JWTSecret string
	JWTTTL    time.Duration
	AppEnv    string
}

func Load() Config {
	return Config{
		HTTPPort:  getenv("HTTP_PORT", "8080"),
		DBDSN:     must("DB_DSN"),     // ej: root:root@tcp(127.0.0.1:3306)/usersdb?parseTime=true
		JWTSecret: must("JWT_SECRET"), // ej: supersecret
		JWTTTL:    getMinutes("JWT_TTL_MIN", 60),
		AppEnv:    getenv("APP_ENV", "local"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}
func getMinutes(k string, def int) time.Duration {
	v := os.Getenv(k)
	if v == "" {
		return time.Duration(def) * time.Minute
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("invalid %s", k)
	}
	return time.Duration(n) * time.Minute
}
