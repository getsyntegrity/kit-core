# AGENTS.md — kit-core

## Authority

This document is the source of truth for how kit-core is designed, implemented, and validated. Any change that conflicts with it is out of scope. When in doubt, preserve the invariants below over feature requests.

---

## Core Principles

- **Determinism** — Same inputs and injected dependencies yield the same outputs. No implicit time or randomness; all variable inputs are explicit and injectable.
- **No hidden side effects** — Functions do not perform I/O, change global state, or depend on environment unless that is their declared purpose and surfaced in their signature or interface.
- **Interface-first design** — Abstractions (Clock, ID, Repository, etc.) are defined as interfaces. Implementations live outside this repo; kit-core owns only the contracts.
- **No global state** — No package-level mutable variables, no singletons. Dependencies are passed in (constructor, handler, or request context).
- **Fail closed** — On ambiguity or error, prefer no behavior over incorrect or non-deterministic behavior. Surface errors; do not hide or default silently.

---

## Architectural Quality Principles

kit-core enforces strict structural discipline. These principles are mandatory.

### High Cohesion

- Each package MUST represent a single, clearly defined responsibility.
- Mixing domain rules with infrastructure or adapter logic inside the same package is FORBIDDEN.
- Package names MUST reflect their responsibility precisely.
- If a package requires unrelated concepts, it MUST be split.

### Low Coupling

- Dependencies MUST point inward only.
- Cross-layer shortcuts are FORBIDDEN.
- Cross-layer imports are FORBIDDEN.
- Circular dependencies are FORBIDDEN.
- Domain MUST NOT depend on infrastructure, adapters, or concrete implementations.
- Application MUST depend only on domain and interfaces.

### Pure Functions First

- Pure functions are first-class citizens.
- Any function that can be pure MUST be pure.
- Business and domain logic MUST prefer pure functions.
- Side effects MUST be isolated to adapters or injected dependencies.
- Hidden mutable state is FORBIDDEN.
- Shared mutable state is FORBIDDEN.
- Package-level mutable variables are FORBIDDEN.

---

## 8. Cohesion and Structural Awareness

Contributors must:

- Recognize when a package violates single responsibility.
- Detect hidden coupling between domain and infrastructure.
- Prefer extracting pure functions over embedding logic in orchestrators.
- Avoid introducing hidden mutable state.
- Understand that determinism is strengthened by purity and cohesion.

---

## Architectural Invariants

- **Domain must be pure** — Logic that expresses business or domain rules must be pure functions of their arguments and injected interfaces. No direct I/O, no `time`, no `math/rand` in that logic.
- **No I/O in domain** — Reading from the network, filesystem, or process environment is not allowed in domain or core library code. I/O happens only in adapters that implement kit-core interfaces.
- **No `time.Now()` in domain** — Time is obtained only via an injected `Clock` (or equivalent) interface. No use of `time.Now()` or similar in code that is part of kit-core’s domain or core types.
- **No randomness in domain** — Identity and any random-like values are produced only via injected abstractions (e.g. ID generator interface). No direct use of `math/rand`, `crypto/rand`, or UUID libraries in domain logic.
- **Clock and ID must be injected via interfaces** — All time and identity sources are dependencies provided by the caller. No default implementations that read from the real clock or system RNG inside kit-core.

---

## Testing Strategy

The project uses **BDD + TDD** with [go-specs](https://github.com/pablogore/go-specs).

### Test philosophy

- **Behavior Driven Development** — Structure tests as describe → when → it. One spec entrypoint per package.
- **Deterministic unit tests** — No flaky tests; no dependence on wall clock, random seed, or environment.
- **Explicit test context** — Time and identity are injected via interfaces; use fakes in tests.
- **No hidden I/O** — Unit tests do not perform network, filesystem, or environment reads.
- **No global state** — Test fixtures and fakes are passed in or created per spec.
- **No reliance on environment variables** — All inputs are set in the test.

### Spec structure

Tests must follow **describe → when → it** using go-specs:

```go
specs.Describe(t, "PackageName", func(s *specs.Spec) {
    s.When("condition or subject", func(s *specs.Spec) {
        s.It("expected behavior", func(ctx *specs.Context) {
            assert.NoError(ctx.T, err)
        })
    })
})
```

- **One entrypoint per package:** `func TestPackageName(t *testing.T) { specs.Describe(t, "packagename", ...) }`
- Use `ctx.T` for testify assertions: `assert.Equal(ctx.T, expected, actual)`

### Test layers

- **Domain tests** — Pure logic; no infrastructure. Use in-memory fakes and fixtures only.
- **Application tests** — Application services with injected domain and adapters; use fakes for repositories, clock, ID generator.
- **Adapter tests** — Adapters implementing kit-core interfaces; may use test doubles from `testkit`.

Domain must not depend on infrastructure. Keep domain specs free of I/O and real time.

### TDD workflow

- Write a **failing spec** first.
- Implement **minimal code** to pass.
- **Refactor**; keep specs green.
- All new code must be written **test-first** when feasible.

### Assertions

Use only:

- `github.com/stretchr/testify/assert`
- `github.com/stretchr/testify/require`

Do not introduce other assertion libraries.

### Deterministic testing

Tests must **not**:

- use `time.Sleep`
- read environment variables
- perform network calls
- access the filesystem
- depend on real clocks or system RNG

When time or identity is needed, inject a **fake clock** or **deterministic ID generator** (e.g. from `testkit`).

### Testkit

Reusable test utilities live under **`testkit`** (root package):

- **fake/** — Deterministic in-memory implementations of interfaces (e.g. clock, repository).
- **spy/** — Record calls and parameters for verification.
- **fixtures/** — Reusable test data (e.g. valid ULIDs, tenant IDs).
- **builders/** — Fluent builders for test entities where helpful.

Fakes and fixtures must be deterministic and must not perform external calls.

---

## Engineering Discipline

- **Deterministic unit tests only** — Tests must be fully deterministic. No flaky tests; no tests that depend on wall clock, random seed, or environment. Use fake clocks and deterministic ID generators in tests.
- **No environment reads in unit tests** — Unit tests must not read from `os.Environ`, config files, or the host. All inputs are set in the test.
- **No network in unit tests** — Unit tests do not open sockets, make HTTP calls, or connect to any external service. Integration or E2E tests that need network belong in a separate suite and are not required for kit-core’s minimal scope.

---

## AI Agent Behavior Expectations

- **Explain reasoning before code** — Before proposing or editing code, state how it fits the principles and invariants above. If a change weakens determinism or introduces I/O in domain, do not propose it.
- **Refuse non-deterministic proposals** — Reject suggestions that add `time.Now()`, `rand`, or environment-dependent behavior to domain or core types. Propose injection via interfaces instead.
- **Stop and ask if ambiguity exists** — If a request could be implemented in a way that breaks invariants (e.g. adding a default implementation that uses real time or network), do not guess. Ask for clarification and insist on interface injection or moving the implementation out of kit-core.

---

## Validation

Invariants and principles are checked by a validation procedure. See **[AGENTS.validation.md](./AGENTS.validation.md)** for the exact steps and criteria.

- [ ] No circular dependencies.
- [ ] No cross-layer imports.
- [ ] Domain contains only pure logic.
- [ ] No package-level mutable variables.
- [ ] Cohesion respected (no mixed concerns in package).
