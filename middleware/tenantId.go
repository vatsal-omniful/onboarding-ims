package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TenantId struct {
	TenantId string `json:"tenant_id" bson:"tenant_id" bind:"required"`
}

func TenantIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tenantId TenantId
		if err := ctx.ShouldBindJSON(&tenantId); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
			ctx.Abort()
			return
		}
		if tenantId.TenantId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "tenant ID is required"})
			ctx.Abort()
			return
		}
		ctx.Set("tenantId", tenantId.TenantId)
		ctx.Next()
	}
}
