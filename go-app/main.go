package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client

	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "handler", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "go_http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "handler", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

type Response struct {
	Source   string  `json:"source,omitempty"`
	Time     string  `json:"time,omitempty"`
	Message  string  `json:"message,omitempty"`
	CacheKey string  `json:"cache_key,omitempty"`
	TTL      float64 `json:"ttl,omitempty"`
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

// customResponseWriter para capturar o status code
type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *customResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func metricsMiddleware(handler string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		crw := &customResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next(crw, r)
		
		duration := time.Since(start).Seconds()
		statusStr := strconv.Itoa(crw.statusCode)
		
		httpRequestsTotal.WithLabelValues(r.Method, handler, statusStr).Inc()
		httpRequestDuration.WithLabelValues(r.Method, handler, statusStr).Observe(duration)
	}
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
		ttlDuration, _ := redisClient.TTL(ctx, cacheKey).Result()
		ttl := ttlDuration.Seconds()
		if ttl < 0 {
			ttl = 0
		}
		json.NewEncoder(w).Encode(Response{Source: "cache", Time: val, CacheKey: cacheKey, TTL: ttl})
		return
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	// Cache por 1 minuto (60 segundos)
	redisClient.Set(ctx, cacheKey, currentTime, 1*time.Minute)

	json.NewEncoder(w).Encode(Response{Source: "server", Time: currentTime, CacheKey: cacheKey, TTL: 60})
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

	http.HandleFunc("/", metricsMiddleware("/", rootHandler))
	http.HandleFunc("/hora", metricsMiddleware("/hora", timeHandler))
	http.HandleFunc("/health", metricsMiddleware("/health", healthHandler))
	http.Handle("/metrics", promhttp.Handler())

	port := "8080"
	fmt.Printf("Servidor Go rodando na porta %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
