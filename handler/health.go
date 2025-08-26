package handler

import "net/http"

// Healthz возвращает "ok" для healthcheck
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
