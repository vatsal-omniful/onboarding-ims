package middleware

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/log"
)

// RequestBodyLoggerMiddleware logs the full JSON body of incoming requests
func RequestBodyLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Only log for POST, PUT, PATCH requests that might have a body
		if ctx.Request.Method == "POST" || ctx.Request.Method == "PUT" || ctx.Request.Method == "PATCH" {
			// Read the request body
			body, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				log.Errorf("Error reading request body: %v", err)
				ctx.Next()
				return
			}

			// Log the full request body with more visibility
			fmt.Printf("=== REQUEST BODY LOG ===\n")
			fmt.Printf("Method: %s\n", ctx.Request.Method)
			fmt.Printf("Path: %s\n", ctx.Request.URL.Path)
			fmt.Printf("Body: %s\n", string(body))
			fmt.Printf("=== END REQUEST BODY ===\n")

			// Also log using the structured logger
			log.Infof("Request Body for %s %s: %s", ctx.Request.Method, ctx.Request.URL.Path, string(body))

			// Restore the body for subsequent middleware/handlers
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		ctx.Next()
	}
}

// DetailedRequestLoggerMiddleware logs detailed request information including headers
func DetailedRequestLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Infof("=== Incoming Request ===")
		log.Infof("Method: %s", ctx.Request.Method)
		log.Infof("URL: %s", ctx.Request.URL.String())
		log.Infof("Content-Type: %s", ctx.GetHeader("Content-Type"))
		log.Infof("User-Agent: %s", ctx.GetHeader("User-Agent"))

		// Log request body for POST, PUT, PATCH requests
		if ctx.Request.Method == "POST" || ctx.Request.Method == "PUT" || ctx.Request.Method == "PATCH" {
			body, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				log.Errorf("Error reading request body: %v", err)
			} else {
				log.Infof("Request Body: %s", string(body))
				// Restore the body
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		}

		log.Infof("=== End Request ===")
		ctx.Next()
	}
}
