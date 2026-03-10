// Package repository provides types and utilities for repository operations,
// including cursor-based pagination and query options.
package repository

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
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
		sortFields = []SortField{{Field: "created_at", Direction: "desc"}}
	}
	return &CursorPagination{Limit: limit, SortFields: sortFields}
}

// WithCursor sets the cursor for pagination.
func (cp *CursorPagination) WithCursor(cursor *Cursor) *CursorPagination {
	cp.Cursor = cursor
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
