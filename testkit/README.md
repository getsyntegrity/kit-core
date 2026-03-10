# testkit

Reusable testing utilities for kit-core specs. See **AGENTS.md** (Testing Strategy) for philosophy and usage.

- **fake/** — Deterministic in-memory implementations (e.g. `Clock` for `clock.Clock`).
- **spy/** — Call recorders for verifying interactions.
- **fixtures/** — Fixed test data (ULIDs, tenant IDs) for deterministic specs.
- **builders/** — Optional fluent builders for test entities; add as needed.

All components must remain deterministic and must not perform external I/O.

**Import paths:**

- `github.com/getsyntegrity/kit-core/testkit/fake`
- `github.com/getsyntegrity/kit-core/testkit/fixtures`
- `github.com/getsyntegrity/kit-core/testkit/spy`
