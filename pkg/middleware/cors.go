package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(corsAllowOrigin string, corsHeaders string) gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		r := c.Request

		origin := r.Header.Get("Origin")
		allowedOrigins := strings.Split(corsAllowOrigin, ",")
		allowAllOrigins := corsAllowOrigin == "*"

		if allowAllOrigins || contains(origin, allowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)

			if requestMethod := r.Method; requestMethod == http.MethodOptions {
				// Handle preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
				allowed := strings.Split(corsHeaders, ",")
				if corsHeaders == "*" {
					// Set default allowed headers
					allowed = []string{
						"Accept",
						"Accept-Encoding",
						"Accept-Language",
						"Authorization",
						"Cache-Control",
						"Connection",
						"Content-Length",
						"Content-Type",
						"Cookie",
						"Host",
						"If-Modified-Since",
						"If-None-Match",
						"Referer",
						"User-Agent",
						// Additional headers from supertokens.GetAllCORSHeaders() can be added here
					}
				}
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowed, ", "))
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Add allowed headers for actual requests
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedOrigins, ", "))
		}

		c.Next()
	}
}

func contains(str string, list []string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}
