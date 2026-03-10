# Prompt: Migrar tests unitarios a go-specs (sclevine/spec)

Usa este prompt con un agente de código (Cursor, etc.) para migrar todos los tests unitarios del repositorio **kit-core** a un estilo BDD con **github.com/sclevine/spec**.

---

## Objetivo

Migrar todos los archivos `*_test.go` del repositorio para usar **sclevine/spec** como organizador BDD, manteniendo:

- **testify** (assert/require) para aserciones — spec no provee assertions; se sigue usando `assert.*` y `require.*`.
- Misma cobertura y casos de prueba; solo cambia la **estructura** (describe/when/it) y, donde aplique, `it.Before` / `it.After`.
- Respeto a **AGENTS.md**: tests deterministas, sin I/O oculto, sin `time.Sleep` en tests que deban ser deterministas, sin lectura de entorno en unit tests.

---

## Librería a usar

- **Spec:** `github.com/sclevine/spec`
- **Añadir dependencia:** `go get github.com/sclevine/spec`
- **Documentación:** [pkg.go.dev/github.com/sclevine/spec](https://pkg.go.dev/github.com/sclevine/spec)

Spec solo organiza subtests (describe/when/it); no reemplaza `testing` ni aporta assertions. Se usa con el `*testing.T` estándar y con testify.

---

## Patrón de migración

### Antes (test estándar + testify)

```go
package idgen

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestIsValidULID(t *testing.T) {
	assert.False(t, IsValidULID(""))
	assert.True(t, IsValidULID("01ARZ3NDEKTSV4RRFFQ69G5FAV"))
	assert.False(t, IsValidULID("invalid"))
}

func TestParseULID(t *testing.T) {
	u, err := ParseULID("01ARZ3NDEKTSV4RRFFQ69G5FAV")
	assert.NoError(t, err)
	assert.Equal(t, "01ARZ3NDEKTSV4RRFFQ69G5FAV", u.String())
	_, err = ParseULID("bad")
	assert.Error(t, err)
}
```

### Después (spec + testify)

```go
package idgen

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/assert"
)

func TestIDGen(t *testing.T) {
	spec.Run(t, "idgen", func(t *testing.T, when spec.G, it spec.S) {
		when("IsValidULID", func() {
			it("returns false for empty string", func() {
				assert.False(t, IsValidULID(""))
			})
			it("returns true for valid ULID", func() {
				assert.True(t, IsValidULID("01ARZ3NDEKTSV4RRFFQ69G5FAV"))
			})
			it("returns false for invalid string", func() {
				assert.False(t, IsValidULID("invalid"))
			})
		})

		when("ParseULID", func() {
			it("parses valid ULID and returns same string", func() {
				u, err := ParseULID("01ARZ3NDEKTSV4RRFFQ69G5FAV")
				assert.NoError(t, err)
				assert.Equal(t, "01ARZ3NDEKTSV4RRFFQ69G5FAV", u.String())
			})
			it("returns error for invalid input", func() {
				_, err := ParseULID("bad")
				assert.Error(t, err)
			})
		})
	})
}
```

---

## Reglas de migración

1. **Un solo `spec.Run` por archivo** (o uno por “suite” si se usa `spec.Suite` en init). El primer argumento de `spec.Run` es un nombre de suite (p. ej. el nombre del paquete o del tipo bajo test).

2. **Nomenclatura BDD:**
   - `when("Algo", func() { ... })` para agrupar por funcionalidad o contexto.
   - `it("descripción en presente, resultado esperado", func() { ... })` para cada caso.
   - Descripciones en **inglés**, claras y legibles (ej.: "returns false for empty string").

3. **Mantener testify:** usar `assert.*` y `require.*` dentro de los bloques `it(...)`; no introducir otra librería de assertions.

4. **Before/After solo si aportan:** usar `it.Before(func() { ... })` e `it.After(func() { ... })` cuando eviten duplicación o configuren estado común; no obligatorio en todos los archivos.

5. **Determinismo (AGENTS.md):** no añadir `time.Sleep`, no leer `os.Environ`, no usar `time.Now()` para lógica de test salvo que ya exista (p. ej. reloj inyectado). Los tests que hoy usan `time.Sleep` (timeout_test, retry_test, resilient_handlers_test) no deben empeorar; idealmente se mantienen o se refactorizan después con reloj falso.

6. **Tablas (table-driven):** si un test es table-driven con muchos casos, se puede:
   - mantener un solo `it("cubre todos los casos de la tabla", func() { ... })` que recorra la tabla, o
   - convertir cada fila en un `it("descripción por fila", func() { ... })` cuando mejore la legibilidad.

7. **Panics:** para tests que comprueban panic, seguir usando `assert.Panics(t, func() { ... })` dentro del bloque `it(...)`.

---

## Archivos a migrar (orden sugerido)

Migrar todos los siguientes `*_test.go`:

1. `idgen/generator_test.go`
2. `validation/errors_test.go`
3. `validation/rules_test.go`
4. `validation/validator_test.go`
5. `timestamp/timestamp_service_test.go`
6. `listeners/listeners_test.go`
7. `fflags/evaluator_test.go`
8. `repository/cursor_pagination_test.go`
9. `repository/query_options_test.go`
10. `repository/queryable_repository_test.go`
11. `repository/repository_test.go`
12. `domain/aggregate_test.go`
13. `domain/errors_test.go`
14. `strategy/query_strategy_test.go`
15. `tenant/context_test.go`
16. `errorschain/errorschain_test.go`
17. `security/config/config_test.go`
18. `infra/capabilities/policy_test.go`

Además, en `pkg/resilience/` (si existen en el árbol actual):

- `pkg/resilience/*_test.go` (config_test, fallback_test, circuit_breaker_test, retry_test, timeout_test, resilient_handlers_test, etc.)

---

## Verificación final

Tras la migración:

1. `go build ./...` debe pasar.
2. `go test ./... -count=1` debe pasar (sin flakiness).
3. No debe quedar ningún `func TestXxx(t *testing.T)` que no esté dentro de un `spec.Run` como punto de entrada; el único `Test*` expuesto al test runner puede ser, por ejemplo, `func TestPackageName(t *testing.T) { spec.Run(t, "packageName", testPackageName) }`.
4. Opcional: `go test -v ./...` para revisar que los nombres de `when`/`it` se vean bien en la salida.

---

## Ejemplo de salida deseada (`go test -v`)

```
=== RUN   TestIDGen
=== RUN   TestIDGen/idgen
=== RUN   TestIDGen/idgen/IsValidULID
=== RUN   TestIDGen/idgen/IsValidULID/returns_false_for_empty_string
=== RUN   TestIDGen/idgen/IsValidULID/returns_true_for_valid_ULID
...
```

---

## Resumen para el agente

- **Añadir:** `github.com/sclevine/spec` al módulo.
- **Refactorizar:** cada `*_test.go` a un único `spec.Run` (o suite) con `when`/`it`, manteniendo `assert`/`require` de testify.
- **No cambiar:** la API del código de producción ni el comportamiento de los tests; solo la organización y los nombres de los casos.
- **Cumplir:** AGENTS.md (determinismo, sin I/O oculto en domain, sin global state en tests).
