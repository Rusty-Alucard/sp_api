package domain

import (
	"reflect"
)

// ValueObject represents a domain value object
// Value objects are immutable and defined by their attributes
type ValueObject interface {
	Equals(other ValueObject) bool
}

// BaseValueObject provides common functionality for value objects
type BaseValueObject struct{}

// Equals compares two value objects for equality
func (v BaseValueObject) Equals(other ValueObject) bool {
	if other == nil {
		return false
	}
	return reflect.DeepEqual(v, other)
}

// StringValueObject represents a string-based value object
type StringValueObject struct {
	BaseValueObject
	value string
}

// NewStringValueObject creates a new string value object
func NewStringValueObject(value string) StringValueObject {
	return StringValueObject{value: value}
}

// Value returns the string value
func (v StringValueObject) Value() string {
	return v.value
}

// String returns the string representation
func (v StringValueObject) String() string {
	return v.value
}

// Equals compares with another value object
func (v StringValueObject) Equals(other ValueObject) bool {
	if other == nil {
		return false
	}
	if otherVO, ok := other.(StringValueObject); ok {
		return v.value == otherVO.value
	}
	return false
}
