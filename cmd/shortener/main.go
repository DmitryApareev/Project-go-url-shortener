// @title URL Shortener API
// @version 1.0
// @description Simple URL shortener
// @BasePath /

package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/DmitryApareev/Project-go-url-shortener/handler"
	"github.com/DmitryApareev/Project-go-url-shortener/internal/http/middleware"
	"github.com/DmitryApareev/Project-go-url-shortener/internal/log"
	"github.com/DmitryApareev/Project-go-url-shortener/store"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/DmitryApareev/Project-go-url-shortener/docs"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

func main() {
	logger, _ := log.New(getenv("LOG_LEVEL", "info"))
	defer logger.Sync()

	r := gin.Default()

	// Healthcheck
	r.GET("/healthz", func(c *gin.Context) {
		handler.Healthz(c.Writer, c.Request)
	})

	// Prometheus metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger UI
	r.GET("/swagger/*any", gin.WrapH(httpSwagger.WrapHandler))

	store.InitializeStore()

	// Rate Limit middleware
	rps := parseFloat(getenv("RATE_LIMIT_RPS", "5"))
	burst := parseInt(getenv("RATE_LIMIT_BURST", "10"))
	r.Use(GinRateLimit(rps, burst))

	// Бизнес-хендлеры
	r.POST("/api/shorten", handler.CreateShortUrl)
	r.GET("/api/:shortUrl", handler.HandleShortUrlRedirect)

	// Server
	srv := &http.Server{
		Addr:         ":" + getenv("PORT", "8080"),
		Handler:      middleware.Logging(logger)(r),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("server error", zap.Any("error", err))
		}
	}()

	logger.Info("server started", zap.String("port", getenv("PORT", "8080")))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Any("error", err))
	}

	logger.Info("server exited properly")
}

// helpers
func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func parseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

// GinRateLimit обёртка для middleware.RateLimit
func GinRateLimit(rps float64, burst int) gin.HandlerFunc {
	limiter := middleware.RateLimit(rps, burst)
	return func(c *gin.Context) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		})
		limiter(handler).ServeHTTP(c.Writer, c.Request)
	}
}
