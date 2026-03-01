package tenant

import (
	"context"

	kitlog "libs/kit-logger/pkg/logger"
)

// ConfigureLoggerCorrelationExtractor configures the logger to extract correlation_id and causation_id from tenant context.
func ConfigureLoggerCorrelationExtractor() {
	kitlog.SetContextFieldExtractor(func(ctx context.Context) []any {
		service := NewTenantContextService()
		tenantCtx, _ := service.GetTenantContext(ctx)
		var fields []any
		if tenantCtx.CorrelationID != "" {
			fields = append(fields, "correlation_id", tenantCtx.CorrelationID)
		}
		if tenantCtx.CausationID != "" {
			fields = append(fields, "causation_id", tenantCtx.CausationID)
		}
		return fields
	})
}
