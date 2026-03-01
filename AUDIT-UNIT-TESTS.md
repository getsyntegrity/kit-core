# Static audit: Go unit tests — external resources and side effects

**Scope:** All `*_test.go` files in the repository.  
**Goal:** Detect unit tests that access real external resources or use forbidden patterns.  
**No code was modified; analysis only.**

---

## 1. Summary

| Metric | Count |
|--------|--------|
| **Total test files** | 9 |
| **Pure unit tests (safe)** | 5 |
| **Integration-style tests** | 0 |
| **Flaky tests (time-dependent)** | 3 |
| **Unsafe tests (external infra/env)** | 1 |

---

## 2. Detailed table

| File | Classification | Reason | Resource / pattern |
|------|----------------|--------|---------------------|
| `pkg/ref/ref_test.go` | **A) Pure** | Only testing, require; no I/O, no time.Sleep, no env, no network. | — |
| `pkg/pii/redact_test.go` | **A) Pure** | Table-driven tests and benchmark; no I/O, no time, no env. | — |
| `pkg/resilience/fallback_test.go` | **A) Pure** | context + errors + testing only; no Sleep, no network, no time.Now(). | — |
| `pkg/resilience/config_test.go` | **A) Pure** | Uses `time` only for constants (e.g. `time.Second`, `time.Millisecond`); no Sleep, no time.Now(). | — |
| `pkg/resilience/circuit_breaker_test.go` | **A) Pure** | No Sleep; uses `FixedClock{Time: time.Now()}` so clock is fixed for the test run. No network/I/O. | `time` (for FixedClock seed only) |
| `pkg/resilience/timeout_test.go` | **C) Flaky** | Uses `time.Sleep` to simulate slow work so timeout triggers. | `time.Sleep` (7 occurrences: 50ms each) |
| `pkg/resilience/retry_test.go` | **C) Flaky** | Uses `time.Sleep(50*time.Millisecond)` in a goroutine to cancel context after a delay (timing sync). | `time.Sleep` (1 occurrence) |
| `pkg/resilience/resilient_handlers_test.go` | **C) Flaky** | Uses `time.Sleep(200*time.Millisecond)` to simulate long-running handler so timeout is hit. | `time.Sleep` (1 occurrence); also `time.Now()` for FixedClock |
| `infra/capabilities/policy_test.go` | **D) Unsafe** | Uses `go/build` and `build.Import(...)` to load package metadata. Reads from Go build context (module cache / GOPATH). | `go/build`; filesystem / env (GOPATH, GOMOD) |

---

## 3. STEP 1 — Forbidden imports (findings)

- **`time`** — Present in: `timeout_test.go`, `retry_test.go`, `resilient_handlers_test.go`, `config_test.go`, `circuit_breaker_test.go`.  
  - **Flagged only where used with `time.Sleep`:** `timeout_test.go`, `retry_test.go`, `resilient_handlers_test.go`.  
  - `config_test.go` and `circuit_breaker_test.go` use `time` only for constants or `FixedClock{Time: time.Now()}` (no Sleep).
- **`go/build`** — Present in `infra/capabilities/policy_test.go`; used for `build.Import` (reads package graph / env).
- **No other forbidden imports** found: no `net`, `net/http`, `google.golang.org/grpc`, `github.com/jackc/pgx`, `go-redis`, `testcontainers`, kafka, franz-go, `go.temporal.io`, `os` (Getenv), or `database/sql` in any `*_test.go`.

---

## 4. STEP 2 — Runtime side effects (findings)

| Pattern | File(s) |
|---------|--------|
| `time.Sleep` | `pkg/resilience/timeout_test.go` (7×), `pkg/resilience/retry_test.go` (1×), `pkg/resilience/resilient_handlers_test.go` (1×) |
| `build.Import` (reads build/env) | `infra/capabilities/policy_test.go` |

**Not found in any test file:**  
`sql.Open`, `pgx.Connect`, `redis.NewClient`, `http.NewRequest`, `http.DefaultClient`, `grpc.Dial`, `net.Dial`, `os.Getenv`, `os.Open`, `exec.Command`, `testcontainers.Run`, docker, `localhost:`, `127.0.0.1:`, or real `http://`/`https://` URLs.

---

## 5. Recommendations

### 5.1 Move to integration suite or refactor (C — Flaky)

- **`pkg/resilience/timeout_test.go`**  
  - **Recommendation:** Refactor to use an injected/fake clock so timeout behavior is asserted without real wall-clock delay.  
  - Alternatively, move the Sleep-based cases to an integration or “timing” test suite and keep unit tests purely clock-fake based.

- **`pkg/resilience/retry_test.go`**  
  - **Recommendation:** Remove the goroutine that uses `time.Sleep(50*time.Millisecond)` then `cancel()`. Use a fake clock or a synchronisation mechanism that does not depend on wall-clock (e.g. inject a clock and advance it in test, or trigger cancel via a channel/callback).

- **`pkg/resilience/resilient_handlers_test.go`**  
  - **Recommendation:** Replace the 200ms `time.Sleep` in the timeout test with a fake clock or a mock that signals “long running” without sleeping (e.g. block on a channel until context is done).

### 5.2 Use mocks / fakes

- **Resilience package (timeout, retry, resilient_handlers):**  
  - Introduce a **Clock** (or similar) interface used by timeout/retry/resilient handler logic and inject it in tests.  
  - In unit tests, use a **fixed or fake clock** (e.g. fixed timestamp or step-advancing clock) so that:  
    - No `time.Sleep` is required.  
    - No dependence on real time.  
  - Aligns with AGENTS.md: deterministic tests, no implicit time, inject time via interfaces.

- **`pkg/resilience/circuit_breaker_test.go` and `resilient_handlers_test.go`:**  
  - Already use `FixedClock{Time: time.Now()}`. For full determinism, seed with a **fixed timestamp** (e.g. `time.Date(...)`) instead of `time.Now()`.

### 5.3 Unsafe test — environment / build (D)

- **`infra/capabilities/policy_test.go`**  
  - **Recommendation:**  
    - Prefer a **static check** (e.g. script or analyzer) that inspects source/imports without loading the full build context, or  
    - Move this to a **build-time / CI check** that is explicitly allowed to depend on the Go environment and module graph, and document it as non-unit (e.g. “build guardrail” or “integration-style guardrail”).  
  - Do not treat it as a pure unit test: it depends on filesystem and Go env (GOPATH/module cache).

### 5.4 In-memory adapters

- No tests in this repo currently connect to real DBs, Redis, Kafka, or gRPC.  
- If future code adds such dependencies, unit tests should use in-memory or fake implementations and no real network or external services.

---

## 6. Summary by classification

- **A) Pure (5 files):** `pkg/ref/ref_test.go`, `pkg/pii/redact_test.go`, `pkg/resilience/fallback_test.go`, `pkg/resilience/config_test.go`, `pkg/resilience/circuit_breaker_test.go`.
- **B) Integration-style:** None.
- **C) Flaky (3 files):** `pkg/resilience/timeout_test.go`, `pkg/resilience/retry_test.go`, `pkg/resilience/resilient_handlers_test.go` — all due to `time.Sleep` for timing.
- **D) Unsafe (1 file):** `infra/capabilities/policy_test.go` — due to `build.Import` and dependency on Go build environment / filesystem.
