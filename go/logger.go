package openapi

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s %s",
			removePort(r.RemoteAddr),
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func removePort(ipAddress string) string {
	colonIndex := strings.LastIndex(ipAddress, ":")
	if colonIndex > -1 {
		return ipAddress[:colonIndex] // Slice the string up to the colon
	}
	return ipAddress
}
