package tenant

// ConfigureLoggerCorrelationExtractor configures the logger to extract correlation_id and causation_id from tenant context.
// Kit-core does not depend on a logger; implementations in runtime may wire this.
func ConfigureLoggerCorrelationExtractor() {
	// No-op: kit-core does not depend on kit-logger. Callers in runtime can wire correlation to their logger.
}
