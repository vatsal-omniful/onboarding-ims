package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TenantId struct {
	TenantId string `json:"tenantId" bson:"tenantId" bind:"required"`
}

func TenantIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Read the request body without consuming it
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error reading request body"})
			ctx.Abort()
			return
		}

		// Restore the body for subsequent handlers
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Parse the JSON to extract tenantId
		var requestData map[string]any
		if err := json.Unmarshal(body, &requestData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
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
