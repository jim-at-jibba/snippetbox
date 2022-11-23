package main

import (
	"net/http"
)

// Becaise we want this to run on ALL requests we must execute it
// BEFORE the servemuc
// secureHeaders -> servemux -> application handler
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XXS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}
