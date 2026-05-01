package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
)

type Response struct {
	Source  string `json:"source,omitempty"`
	Time    string `json:"time,omitempty"`
	Message string `json:"message,omitempty"`
}

func initRedis() {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
		DB:   0,
	})
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "Olá! Esta é a aplicação em Go respondendo com um texto fixo."})
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cacheKey := "go_app_time"

	val, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		json.NewEncoder(w).Encode(Response{Source: "cache", Time: val})
		return
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	// Cache por 1 minuto (60 segundos)
	redisClient.Set(ctx, cacheKey, currentTime, 1*time.Minute)

	json.NewEncoder(w).Encode(Response{Source: "server", Time: currentTime})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	status := "up"
	redisStatus := "ok"
	
	if err := redisClient.Ping(ctx).Err(); err != nil {
		status = "degraded"
		redisStatus = "error"
	}
	
	json.NewEncoder(w).Encode(map[string]string{
		"status": status,
		"redis":  redisStatus,
	})
}

func main() {
	initRedis()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/hora", timeHandler)
	http.HandleFunc("/health", healthHandler)
	http.Handle("/metrics", promhttp.Handler())

	port := "8080"
	fmt.Printf("Servidor Go rodando na porta %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
