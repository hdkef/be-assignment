package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
)

// VerifySessionMiddleware adapts the Supertoken VerifySession function to work as a Gin middleware
func VerifySessionMiddleware(options *sessmodels.VerifySessionOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sessionVerified bool
		var err error

		session.VerifySession(options, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Session verified successfully
			sessionVerified = true
			// Update the Gin context with the new request that has the session information
			c.Request = r
		})).ServeHTTP(c.Writer, c.Request)

		if !sessionVerified {
			// Session verification failed
			err = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			if err != nil {
				// If there's an error sending the response, log it
				fmt.Printf("Error sending unauthorized response: %v\n", err)
			}
			return
		}

		// Continue with the next middleware/handler
		c.Next()
	}
}
