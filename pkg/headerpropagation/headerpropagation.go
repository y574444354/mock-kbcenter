package headerpropagation

import (
	"context"

	"github.com/zgsm/mock-kbcenter/config"
)

type ContextKey string

const (
	// ContextKeyPrefix is the prefix for context keys
	ContextKeyPrefix ContextKey = "header."
)

// GetHeaderValue gets propagated header value from context
func GetHeaderValue(ctx context.Context, header string) string {
	if val := ctx.Value(ContextKeyPrefix + ContextKey(header)); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetAllPropagatedHeaders gets all propagated headers from context
func GetAllPropagatedHeaders(ctx context.Context) map[string]string {
	headers := make(map[string]string)
	if cfg := config.GetConfig(); cfg != nil {
		for _, header := range cfg.HeaderPropagation.Headers {
			if val := GetHeaderValue(ctx, header); val != "" {
				headers[header] = val
			}
		}
	}
	return headers
}

// WithContext injects headers into a new context
func WithContext(ctx context.Context, headers map[string]string) context.Context {
	for k, v := range headers {
		ctx = context.WithValue(ctx, ContextKeyPrefix+ContextKey(k), v)
	}
	return ctx
}
