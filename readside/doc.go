// Package readside provides read-side projection infrastructure.
//
// This package provides offset store implementations for checkpoint persistence.
// It abstracts database operations for storing and retrieving projection checkpoints.
//
// Responsibilities:
//   - Checkpoint persistence: Load and save projection versions
//   - Version contiguity: Enforce contiguous version saves
//   - Database abstraction: Querier interface for database operations
//
// Non-responsibilities:
//   - Event processing (handled by projection listeners)
//   - Event source implementation (handled by EventStream)
//   - Event adaptation (handled by EventAdapter)
package readside
