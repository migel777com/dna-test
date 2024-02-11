package middleware

import (
	"dna-test/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

// Gin set content-type as application/text after AbortWithError and we can not reset it even with c.JSON,
// so I added middleware that sets content-type header at the start
func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

// Can pass the logger or custom logic to errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		var err error
		for _, err = range c.Errors {
			fmt.Errorf("error: %v", err.Error())
		}

		// status -1 doesn't overwrite existing status code
		if err != nil {
			var errAdvanced models.AdvancedErrorResponse
			switch {
			case errors.As(err, &errAdvanced):
				c.JSON(-1, gin.H{"critical": models.AdvancedErrorResponse{
					Message:     errAdvanced.Message,
					Description: errAdvanced.Description,
				}})
			default:
				c.JSON(-1, gin.H{"error": models.ErrorResponse{Message: err.Error()}})
			}
		}
	}
}
