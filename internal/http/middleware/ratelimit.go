package middleware

import (
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"golang.org/x/time/rate"
)

// ipLimiter хранит лимитеры для каждого IP
type ipLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	r        rate.Limit
	burst    int
}

// newIPLimiter создаёт новый лимитер
func newIPLimiter(r rate.Limit, burst int) *ipLimiter {
	return &ipLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		burst:    burst,
	}
}

// get возвращает лимитер для конкретного IP
func (l *ipLimiter) get(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()
	if lim, ok := l.limiters[ip]; ok {
		return lim
	}
	lim := rate.NewLimiter(l.r, l.burst)
	l.limiters[ip] = lim
	return lim
}

// RateLimit middleware ограничивает количество запросов по IP
func RateLimit(rps float64, burst int) func(http.Handler) http.Handler {
	l := newIPLimiter(rate.Limit(rps), burst)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			if !l.get(ip).Allow() {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ====== Тест ======

func TestRateLimit(t *testing.T) {
	// создаём middleware с лимитом 1 запрос, burst 1
	limiter := RateLimit(1, 1)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	// первый запрос должен пройти
	w := httptest.NewRecorder()
	limiter.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatal("first request should pass")
	}

	// второй запрос сразу же — должен заблокироваться
	w2 := httptest.NewRecorder()
	limiter.ServeHTTP(w2, req)
	if w2.Code != 429 {
		t.Fatal("second request should be limited")
	}
}
