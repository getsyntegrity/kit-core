package repository

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Clock provides the current time. Inject in kit-core; no time.Now() in core.
type Clock interface {
	Now() time.Time
}

// Sort direction constants (used by QueryOptions and cursor logic).
const (
	SortDirectionAsc  = "asc"
	SortDirectionDesc = "desc"
)

// Cursor represents a pagination cursor for efficient pagination.
type Cursor struct {
	Direction map[string]string      `json:"direction"`
	Timestamp time.Time              `json:"timestamp"`
	Values    map[string]interface{} `json:"values"`
}

// CursorPagination represents cursor-based pagination options.
type CursorPagination struct {
	Cursor     *Cursor     `json:"cursor,omitempty"`
	Limit      int         `json:"limit"`
	Reverse    bool        `json:"reverse,omitempty"`
	SortFields []SortField `json:"sort_fields"`
}

// SortField represents a field to sort by with its direction.
type SortField struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

// CursorResult represents the result of a cursor-based query.
type CursorResult[T any] struct {
	Data           []T     `json:"data"`
	HasMore        bool    `json:"has_more"`
	HasPrevious    bool    `json:"has_previous"`
	NextCursor     *Cursor `json:"next_cursor,omitempty"`
	PreviousCursor *Cursor `json:"previous_cursor,omitempty"`
	TotalCount     *int64  `json:"total_count,omitempty"`
}

// NewCursorPagination creates a new cursor pagination instance.
func NewCursorPagination(limit int, sortFields ...SortField) *CursorPagination {
	if len(sortFields) == 0 {
		sortFields = []SortField{{Field: "created_at", Direction: SortDirectionDesc}}
	}
	return &CursorPagination{
		Limit:      limit,
		SortFields: sortFields,
	}
}

// WithCursor sets the cursor for pagination.
func (cp *CursorPagination) WithCursor(cursor *Cursor) *CursorPagination {
	cp.Cursor = cursor
	return cp
}

// WithReverse sets the reverse flag for backward pagination.
func (cp *CursorPagination) WithReverse(reverse bool) *CursorPagination {
	cp.Reverse = reverse
	return cp
}

// WithSortFields sets the sort fields.
func (cp *CursorPagination) WithSortFields(sortFields ...SortField) *CursorPagination {
	cp.SortFields = sortFields
	return cp
}

// Validate validates the cursor pagination options.
func (cp *CursorPagination) Validate() error {
	if cp.Limit <= 0 {
		return fmt.Errorf("limit must be positive")
	}
	if cp.Limit > 1000 {
		return fmt.Errorf("limit cannot exceed 1000")
	}
	if len(cp.SortFields) == 0 {
		return fmt.Errorf("at least one sort field must be specified")
	}
	for _, sf := range cp.SortFields {
		if sf.Field == "" {
			return fmt.Errorf("sort field name cannot be empty")
		}
		if sf.Direction != SortDirectionAsc && sf.Direction != SortDirectionDesc {
			return fmt.Errorf("sort direction must be 'asc' or 'desc', got: %s", sf.Direction)
		}
	}
	return nil
}

// BuildCursorWhereClause builds a WHERE clause for cursor-based pagination.
func (cp *CursorPagination) BuildCursorWhereClause() (string, []interface{}, error) {
	if cp.Cursor == nil {
		return "", []interface{}{}, nil
	}
	if err := cp.Validate(); err != nil {
		return "", nil, fmt.Errorf("invalid cursor pagination: %w", err)
	}
	var conditions []string
	var args []interface{}
	argIndex := 1
	for index, sortField := range cp.SortFields {
		cursorValue, exists := cp.Cursor.Values[sortField.Field]
		if !exists {
			continue
		}
		var condition string
		if cp.Reverse {
			if sortField.Direction == SortDirectionAsc {
				condition = fmt.Sprintf("%s < $%d", sortField.Field, argIndex)
			} else {
				condition = fmt.Sprintf("%s > $%d", sortField.Field, argIndex)
			}
		} else {
			if sortField.Direction == SortDirectionAsc {
				condition = fmt.Sprintf("%s > $%d", sortField.Field, argIndex)
			} else {
				condition = fmt.Sprintf("%s < $%d", sortField.Field, argIndex)
			}
		}
		if index > 0 {
			var equalityConditions []string
			for j := 0; j < index; j++ {
				prevSF := cp.SortFields[j]
				prevValue, exists := cp.Cursor.Values[prevSF.Field]
				if exists {
					equalityConditions = append(equalityConditions, fmt.Sprintf("%s = $%d", prevSF.Field, argIndex))
					args = append(args, prevValue)
					argIndex++
				}
			}
			if len(equalityConditions) > 0 {
				condition = fmt.Sprintf("(%s) OR (%s AND %s)",
					strings.Join(equalityConditions, " AND "),
					strings.Join(equalityConditions, " AND "),
					condition)
			}
		}
		conditions = append(conditions, condition)
		args = append(args, cursorValue)
		argIndex++
	}
	if len(conditions) == 0 {
		return "", []interface{}{}, nil
	}
	return fmt.Sprintf("(%s)", strings.Join(conditions, " OR ")), args, nil
}

// BuildOrderByClause builds an ORDER BY clause for cursor pagination.
func (cp *CursorPagination) BuildOrderByClause() string {
	if len(cp.SortFields) == 0 {
		return ""
	}
	var orderClauses []string
	for _, sf := range cp.SortFields {
		orderClauses = append(orderClauses, fmt.Sprintf("%s %s", sf.Field, sf.Direction))
	}
	return strings.Join(orderClauses, ", ")
}

// BuildLimitClause builds the LIMIT clause.
func (cp *CursorPagination) BuildLimitClause() string {
	return fmt.Sprintf("LIMIT %d", cp.Limit+1)
}

// CreateCursorFromEntity creates a cursor from an entity using the provided clock.
func CreateCursorFromEntity[T any](clock Clock, entity T, sortFields []SortField, getFieldValue func(T, string) interface{}) *Cursor {
	if clock == nil || len(sortFields) == 0 {
		return nil
	}
	cursor := &Cursor{
		Values:    make(map[string]interface{}),
		Direction: make(map[string]string),
		Timestamp: clock.Now(),
	}
	for _, sf := range sortFields {
		value := getFieldValue(entity, sf.Field)
		if value != nil {
			cursor.Values[sf.Field] = value
			cursor.Direction[sf.Field] = sf.Direction
		}
	}
	return cursor
}

// CreateCursorFromValues creates a cursor from raw values using the provided clock.
func CreateCursorFromValues(clock Clock, values map[string]interface{}, sortFields []SortField) *Cursor {
	if clock == nil || len(sortFields) == 0 {
		return nil
	}
	cursor := &Cursor{
		Values:    make(map[string]interface{}),
		Direction: make(map[string]string),
		Timestamp: clock.Now(),
	}
	for _, sf := range sortFields {
		if value, exists := values[sf.Field]; exists {
			cursor.Values[sf.Field] = value
			cursor.Direction[sf.Field] = sf.Direction
		}
	}
	return cursor
}

// EncodeCursor encodes a cursor to a base64 string for transport.
func EncodeCursor(cursor *Cursor) (string, error) {
	if cursor == nil {
		return "", nil
	}
	jsonData, err := json.Marshal(cursor)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cursor: %w", err)
	}
	return base64.URLEncoding.EncodeToString(jsonData), nil
}

// DecodeCursor decodes a cursor from a base64 string.
func DecodeCursor(encoded string) (*Cursor, error) {
	if encoded == "" {
		return nil, nil
	}
	jsonData, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cursor: %w", err)
	}
	var cursor Cursor
	if err := json.Unmarshal(jsonData, &cursor); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cursor: %w", err)
	}
	return &cursor, nil
}

// ProcessCursorResult processes the raw query result to create a proper cursor result.
func ProcessCursorResult[T any](clock Clock, data []T, limit int, sortFields []SortField, getFieldValue func(T, string) interface{}) *CursorResult[T] {
	result := &CursorResult[T]{Data: data}
	if len(data) > limit {
		result.Data = data[:limit]
		result.HasMore = true
		if len(result.Data) > 0 {
			lastItem := result.Data[len(result.Data)-1]
			result.NextCursor = CreateCursorFromEntity(clock, lastItem, sortFields, getFieldValue)
		}
	}
	if len(result.Data) > 0 {
		firstItem := result.Data[0]
		result.PreviousCursor = CreateCursorFromEntity(clock, firstItem, sortFields, getFieldValue)
	}
	return result
}

// CursorPaginationBuilder provides a fluent interface for building cursor pagination.
type CursorPaginationBuilder struct {
	pagination *CursorPagination
}

// NewCursorPaginationBuilder creates a new cursor pagination builder.
func NewCursorPaginationBuilder(limit int) *CursorPaginationBuilder {
	return &CursorPaginationBuilder{pagination: NewCursorPagination(limit)}
}

// WithCursor sets the cursor.
func (b *CursorPaginationBuilder) WithCursor(cursor *Cursor) *CursorPaginationBuilder {
	b.pagination.WithCursor(cursor)
	return b
}

// WithReverse sets the reverse flag.
func (b *CursorPaginationBuilder) WithReverse(reverse bool) *CursorPaginationBuilder {
	b.pagination.WithReverse(reverse)
	return b
}

// WithSortField adds a sort field.
func (b *CursorPaginationBuilder) WithSortField(field, direction string) *CursorPaginationBuilder {
	b.pagination.SortFields = append(b.pagination.SortFields, SortField{Field: field, Direction: direction})
	return b
}

// WithSortFields sets multiple sort fields.
func (b *CursorPaginationBuilder) WithSortFields(sortFields ...SortField) *CursorPaginationBuilder {
	b.pagination.SortFields = sortFields
	return b
}

// Build returns the built cursor pagination.
func (b *CursorPaginationBuilder) Build() *CursorPagination {
	return b.pagination
}
