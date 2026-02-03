package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/igoventura/fintrack-core/domain"
)

const (
	TenantIDHeader = "X-Tenant-ID"
)

type TenantMiddleware struct{}

func NewTenantMiddleware() *TenantMiddleware {
	return &TenantMiddleware{}
}

func (m *TenantMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader(TenantIDHeader)
		if tenantID != "" {
			ctx := domain.WithTenantID(c.Request.Context(), tenantID)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}
