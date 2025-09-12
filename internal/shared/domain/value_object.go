// Package domain provides shared domain concepts and base implementations.
package domain

// ValueObject represents the base interface for value objects.
type ValueObject interface {
	Equals(other ValueObject) bool
}

// BaseValueObject provides a base implementation for value objects.
type BaseValueObject struct{}

// Equals checks if two value objects are equal.
func (v BaseValueObject) Equals(_ ValueObject) bool {
	return false // Override in concrete implementations
}
