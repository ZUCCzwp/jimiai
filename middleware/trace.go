package middleware

import (
	"context"
	"fmt"
	"jiyu/global"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type TraceConfig struct {
	RequestHeaderKey string
	ResponseTraceKey string
}

var DefaultTraceConfig = TraceConfig{
	RequestHeaderKey: "X-Request-Id",
	ResponseTraceKey: "X-Trace-Id",
}

func Trace() gin.HandlerFunc {
	return TraceWithConfig(DefaultTraceConfig)
}

func TraceWithConfig(config TraceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		traceID := c.GetHeader(config.RequestHeaderKey)
		if traceID == "" {
			traceID = fmt.Sprintf("%s", strings.ToUpper(xid.New().String()))
		}

		ctx := NewTraceID(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(config.ResponseTraceKey, traceID)
		c.Next()
	}
}
func NewTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, global.TraceIDKey{}, traceID)
}

func FromTraceID(ctx context.Context) string {
	v := ctx.Value(global.TraceIDKey{})
	if v != nil {
		return v.(string)
	}
	return ""

}
