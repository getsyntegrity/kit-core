// Package listeners provides generic listener infrastructure for event processing.
//
// This package provides the core listener implementation that projection listeners use.
// It handles the event processing loop with retry logic for checkpoint persistence.
//
// Responsibilities:
//   - Event processing loop: receive → handle → save checkpoint
//   - Retry logic: Exponential backoff for checkpoint persistence
//   - Error handling: Distinguishes retryable vs non-retryable errors
//
// Non-responsibilities:
//   - Event source implementation (handled by EventStream)
//   - Event adaptation (handled by EventAdapter)
//   - Handler routing (handled by ProjectionHandlerRouter)
//   - Checkpoint persistence implementation (handled by OffsetStore)
package listeners
