package mw

import (
	"log"
	"net/http"
	"time"
)

func LogHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s\n", r.Method, r.URL)
		next(w, r)
		log.Printf("time: %f\n\n", time.Since(start).Seconds())
	}
}
