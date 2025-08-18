package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/log"
)

type TenantId struct {
	TenantId string `json:"tenantId" bson:"tenantId" bind:"required"`
}

func TenantIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// First check if tenantId is in query parameters (for GET requests)
		if tenantIdQuery := ctx.Query("tenantId"); tenantIdQuery != "" {
			tenantIdUint, err := strconv.ParseUint(tenantIdQuery, 10, 64)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "tenantId must be a valid unsigned integer"})
				ctx.Abort()
				return
			}
			ctx.Set("tenantId", uint(tenantIdUint))
			ctx.Next()
			return
		}

		// Read the request body without consuming it
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error reading request body"})
			ctx.Abort()
			return
		}

		// Restore the body for subsequent handlers
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Handle empty body for GET requests
		if len(body) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "tenant ID is required (provide in query parameter or request body)"})
			ctx.Abort()
			return
		}

		// Parse the JSON to extract tenantId
		var requestData map[string]any
		if err := json.Unmarshal(body, &requestData); err != nil {
			log.Errorf("Error unmarshalling request body: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			ctx.Abort()
			return
		}

		// Check if tenantId exists in the request
		tenantIdValue, exists := requestData["tenantId"]
		if !exists {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "tenant ID is required"})
			ctx.Abort()
			return
		}

		ctx.Set("tenantId", tenantIdValue)
		ctx.Next()
	}
}
