// Package spy provides recorders for verifying interactions in tests.
package spy

import "sync"

// CallRecord holds a single recorded call (method name and arguments).
type CallRecord struct {
	Name string
	Args []interface{}
}

// CallRecorder records method calls for later assertion.
// Thread-safe.
type CallRecorder struct {
	mu     sync.Mutex
	Calls  []CallRecord
	Counts map[string]int
}

// NewCallRecorder returns a new call recorder.
func NewCallRecorder() *CallRecorder {
	return &CallRecorder{
		Calls:  nil,
		Counts: make(map[string]int),
	}
}

// Record records a call with the given name and arguments.
func (r *CallRecorder) Record(name string, args ...interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Calls = append(r.Calls, CallRecord{Name: name, Args: args})
	r.Counts[name]++
}

// CallCount returns how many times the given name was recorded.
func (r *CallRecorder) CallCount(name string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Counts[name]
}

// Reset clears all recorded calls.
func (r *CallRecorder) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Calls = nil
	r.Counts = make(map[string]int)
}
