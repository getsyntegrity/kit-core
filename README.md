# kit-core
pablo

Minimal, universal infrastructure foundation for Go services. Abstractions for time, identity, and persistence that keep application code testable and portable.

---

## Design Philosophy

kit-core is a thin dependency layer that defines interfaces and small building blocks only. It does not pick transports, brokers, or storage backends. Application code depends on these abstractions; concrete implementations live in other packages or in the calling project. The goal is deterministic, platform-neutral core types that any Go backend can reuse without pulling in servers, SDKs, or framework opinions.

---

## What It Provides

- **Clock** — Abstraction over time (now, timestamps). Enables fixed or fake clocks in tests and reproducible runs.
- **ID** — Abstraction over identity generation (e.g. UUIDs, ULIDs, or domain-specific IDs). Keeps ID strategy swappable and testable.
- **Repository** — Minimal persistence interface (e.g. load/save by ID). No query language or storage specifics; adapters implement the interface for the store of your choice.
- Small, stable types and interfaces only. No wiring, no frameworks, no platform lock-in.

---

## What It Explicitly Does NOT Provide

- HTTP or RPC server setup
- Message brokers or queues
- Database drivers, query builders, or migrations
- Security or auth providers
- Observability or tracing wiring
- Business logic or domain rules
- Any concrete implementation of transport, broker, or storage

Implementations belong in separate modules or in the application.

---

## Determinism

Deterministic behavior is central: the same inputs and environment should yield the same outcomes. That is only possible if time and identity are controlled.

- **Clock** — Code that needs “now” or timestamps uses a `Clock` interface. In production you inject a real clock; in tests you inject a fixed or stepped clock so that timestamps and time-based logic are reproducible.
- **ID** — Identity generation is behind an abstraction (e.g. ID generator interface). Production uses your chosen scheme (UUID, ULID, etc.); tests use deterministic or predictable IDs so that snapshots, ordering, and assertions are stable.

No direct use of `time.Now()` or ad-hoc ID generation in code that depends on kit-core abstractions; all such behavior goes through these interfaces.

---

## Dependency Rules

kit-core must remain a leaf package with zero dependencies on:

- Any HTTP or RPC server library
- Any message broker or queue client
- Any database driver or ORM
- Any observability or tracing SDK
- Any security or auth provider
- Any framework that forces transport, storage, or platform choices

If a dependency would pull in servers, brokers, storage, or platform-specific SDKs, it must not be added. kit-core defines contracts only; the dependency direction is always: application and adapters depend on kit-core, never the reverse.

---

## Example Usage

```go
package main

import (
	"context"
	"time"
)

// Clock abstraction: production uses real time; tests use a fixed clock.
type Clock interface {
	Now() time.Time
}

// Repository abstraction: load/save by ID; adapter implements for your store.
type Repository[T any] interface {
	Get(ctx context.Context, id string) (T, error)
	Put(ctx context.Context, id string, v T) error
}

func processOrder(ctx context.Context, r Repository[Order], c Clock) error {
	order, err := r.Get(ctx, "order-1")
	if err != nil {
		return err
	}
	order.ProcessedAt = c.Now()
	return r.Put(ctx, order.ID, order)
}
```

---

## Stability and Scope

- **Stable** — Existing interfaces and types will not be changed in breaking ways. New abstractions may be added when they stay transport-, broker-, and storage-agnostic.
- **Scope** — Only interfaces and minimal types that support determinism, identity, and persistence abstraction. No new categories of functionality (servers, brokers, storage, observability) will be added.

---

## License

[Specify license, e.g. MIT, Apache-2.0.]
