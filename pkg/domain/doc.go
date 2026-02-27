// Package domain provides generic domain patterns and interfaces for event sourcing
// and domain-driven design. This package contains reusable components that can be
// used across different microservices to maintain consistency in domain modeling.
//
// Key Components:
//
//   - Aggregate: Interface and base implementation for domain aggregates
//   - EventStore: Interface for event store operations
//   - Repository: Generic repository interfaces for data access
//   - Errors: Domain-specific error types and utilities
//
// Usage Example:
//
//	type User struct {
//		*domain.BaseAggregate
//		Name  string
//		Email string
//	}
//
//	func NewUser(id, name, email string) *User {
//		user := &User{
//			BaseAggregate: domain.NewBaseAggregate(id),
//			Name:          name,
//			Email:         email,
//		}
//		user.AddEvent(&UserCreated{...})
//		return user
//	}
//
//	func (u *User) applyEvent(event proto.Message) error {
//		switch e := event.(type) {
//		case *UserCreated:
//			u.Name = e.Name
//			u.Email = e.Email
//		}
//		return nil
//	}
package domain
