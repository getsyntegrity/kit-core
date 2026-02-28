/*
Package repository provides data access contracts and utilities.

Layout:
  - contracts.go    — Repository, DomainRepository, BaseRepository, Pagination, Result, DomainEntity
  - queryable.go    — QueryableRepository, BaseQueryableRepository, QueryBuilder, QueryOptionsBuilder, QueryResult
  - query_options.go — QueryOptions and its methods
  - cursor.go       — Cursor, CursorPagination, CursorResult, ProcessCursorResult, Clock (pagination utilities)
  - transaction.go   — UnitOfWork, TransactionOptions (driver-agnostic)
  - errors.go       — Error type and helpers (NotFoundError, DatabaseError, etc.)
  - memory/         — In-memory implementation for testing (MemoryRepository)

Implementations (PostgreSQL, Redis, etc.) live in kit-runtime or other modules.
*/
package repository
