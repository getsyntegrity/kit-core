# kit-core v0.1.0 — Public API Surface

This document is the public contract of kit-core v0.1.0. Structural module only; no runtime infra.

---

## domain

**Events / event store**
- `Aggregate` (interface)
- `BaseAggregate`, `NewBaseAggregate`
- `EventStore` (interface)
- `EventRecord` (struct)
- `EventSerializer` (interface)
- `EventMessage` (interface, embeds proto.Message)
- `EventBus` (interface: Publish, Subscribe, Close — uses watermill message type in signature)

**Commands**
- `Command[T]`, `CommandHandler[T]` (interfaces)
- `CommandBus` (interface)

**Entities and repositories**
- `Entity` (interface: GetID, GetType)
- `BaseRepository[T]`, `ReadRepository[T]`, `WriteRepository[T]` (interfaces)
- `TenantScopedRepository[T]`, `TenantScopedReadRepository[T]`, `TenantScopedWriteRepository[T]` (interfaces)
- `RepositoryFactory`, `RepositoryConfig` (interface, struct)
- `MockableEventMessage`, `MockableBaseRepository`, `MockableReadRepository`, `MockableWriteRepository`, `MockableRepositoryFactory` (interfaces)
- `ApplicationContainer[T]` (interface)

**Errors**
- `Error` (struct), `NewError`, `NewErrorWithCause`
- `ErrEntityNotFound`, `ErrNotFound`, `ErrInvalidState`, `ErrConcurrency`, `ErrValidation`, `ErrUnauthorized`, `ErrForbidden`, `ErrInternal`, `ErrExternalService`
- `ErrCode*` constants
- `IsNotFound`, `IsInvalidState`, `IsConcurrency`, `IsValidation` (funcs)

---

## repository

**Interfaces**
- `Repository[T]`, `DomainRepository[T]` (tenant-scoped CRUD)
- `QueryableRepository[T]` (extends Repository with Query, QueryWithCount, QueryWithCursor, FindByField, FindByFields, Search, GetByIDs, BulkSave, BulkDelete, SoftDelete, Restore)
- `DomainEntity` (interface)
- `UnitOfWork` (interface: Do with TransactionOptions)
- `Clock` (interface: Now)

**Structs and types**
- `Pagination`, `Result[T]`, `BaseRepository[T]`, `NewBaseRepository`, `ValidatePagination`, `BuildPagination`
- `QueryOptions`, `NewQueryOptions`, and all `With*` / `Validate` / `Build*` / `Clone` methods
- `QueryResult[T]`, `NewQueryResult`
- `QueryOptionsBuilder`, `NewQueryOptionsBuilder`, `QueryBuilder[T]`, `NewQueryBuilder`
- `BaseQueryableRepository[T]`, `NewBaseQueryableRepository`
- `TransactionOptions`
- `Cursor`, `CursorPagination`, `CursorResult[T]`, `SortField`, `SortDirectionAsc`/`SortDirectionDesc`
- `NewCursorPagination`, `WithCursor`, `WithReverse`, `WithSortFields`, `Validate`, `BuildCursorWhereClause`, `BuildOrderByClause`, `BuildLimitClause`
- `CreateCursorFromEntity`, `CreateCursorFromValues`, `EncodeCursor`, `DecodeCursor`, `ProcessCursorResult`
- `CursorPaginationBuilder`, `NewCursorPaginationBuilder`, `WithCursor`, `WithReverse`, `WithSortField`, `WithSortFields`, `Build`
- `Error`, `NewError`, `NotFoundError`, `AlreadyExistsError`, `InvalidInputError`, `DatabaseError`, `ConnectionError`
- `ErrCodeNotFound`, `ErrCodeAlreadyExists`, `ErrCodeInvalidInput`, `ErrCodeDatabase`, `ErrCodeConnection`

**Subpackage**
- `repository/memory`: `MemoryRepository` (in-memory implementation for testing)

---

## timestamp (clock interfaces)

**Interfaces**
- `TimestampServiceInterface`
- `TimeProvider` (Now)
- `TimestampFormatter` (Format, FormatWithLayout, Parse, ParseWithLayout)

**Types**
- `MockTimeProvider`, `NewMockTimeProvider`, `SetFixedTime`, `AdvanceTime`
- `DefaultTimestampFormatter`, `NewDefaultTimestampFormatter`, `NewDefaultTimestampFormatterWithLayout`, and methods
- `TimestampService`, `NewTimestampServiceWithProvider` (implementation using injected TimeProvider)

---

## observability (minimal contracts)

- `Field` (struct: Key, Value)
- `Logger` (interface: Info, Error, Debug, Warn, WithContext, With)
- `MetricsRecorder` (interface: Inc, Observe)

---

## Other structural packages (exported surface)

- **strategy**: `QueryStrategy[T]`, `SimpleQueryStrategy`, `PaginatedQueryStrategy`, etc.
- **validation**: validators, rules, `Validator`, errors (no HTTP/server)
- **fflags**: feature flag evaluator interfaces/contracts
- **ref**: reference types
- **pii**: redaction, `SensitiveString`
- **errorschain**: error chaining (multierr)
- **setutils**: `Set`, optional, constants
- **resilience**: `Clock` (interface), retry, timeout, fallback, circuit breaker (config + handlers; implementations use injected Clock)
- **infra/capabilities**: policy/config contracts
- **ref**, **pii**, **errorschain**: supporting types and helpers

---

*Generated for final release hardening. No behavior or API changes.*
