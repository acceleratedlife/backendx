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

		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			// Fall back to r.RemoteAddr if X-Forwarded-For is not present
			ip = removePort(r.RemoteAddr)
		} else {
			// X-Forwarded-For may contain a list of IPs, so take the first one
			ip = strings.Split(ip, ",")[0]
		}

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
