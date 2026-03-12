# DevForge + golangci-lint v2

## Problem

When running `forge pr --profile go-lib`, the static-analysis step fails with:

```
go: github.com/golangci/golangci-lint/v2@v2.1.0: module found, but does not contain package github.com/golangci/golangci-lint/v2
```

DevForge invokes:

```bash
go run github.com/golangci/golangci-lint/v2@v2.1.0 run --timeout=5m
```

In golangci-lint v2, the main package is **not** at the module root; it lives under `cmd/golangci-lint`.

## Fix (in pablogore/devforge)

Change the static-analysis invocation from:

- `go run github.com/golangci/golangci-lint/v2@v2.1.0 run ...`

to:

- `go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.0 run ...`

Alternatively, run the `golangci-lint` binary when it exists in PATH (e.g. after `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.0`).

## Workaround in kit-core

Until DevForge is updated:

1. **Lint locally** (use the correct v2 cmd path):
   ```bash
   make lint
   ```
   or:
   ```bash
   go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.0 run --timeout=5m
   ```

2. Then run `forge pr --profile go-lib`. Forge will still fail on the static-analysis step until DevForge is fixed; the steps before it (tidy, conventional-commit, architectural-guard) will pass.

CI may still fail on the forge static-analysis step; the fix must be applied in the DevForge repository (pablogore/devforge).
