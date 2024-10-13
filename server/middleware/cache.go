package middleware

import (
	"net/http"
	"time"
)

func CacheWhileServerIsRunning(next http.Handler) http.Handler {
	startTime := time.Now().Truncate(time.Second)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ifModifiedSince := r.Header.Get("If-Modified-Since")
		if ifModifiedSince != "" {
			modTime, err := time.Parse(http.TimeFormat, ifModifiedSince)
			if err == nil {
				if !startTime.After(modTime) {
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}
		}
		w.Header().Set("Last-Modified", startTime.UTC().Format(http.TimeFormat))
		next.ServeHTTP(w, r)
	})
}
