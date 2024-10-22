package utils

import (
	"errors"
	"fmt"
)

// Tuple is an interface for tuple types.
type Tuple[A, B any] interface {
	First() A
	Second() B
	Size() int
	Get(i int) (any, error)
	String() string
}

// Pair is a type that holds two values of generic types.
type Pair[A, B any] struct {
	first  A
	second B
}

// First returns the first value of the Pair.
func (p Pair[A, B]) First() A {
	return p.first
}

// Second returns the second value of the Pair.
func (p Pair[A, B]) Second() B {
	return p.second
}

// Size returns the size of the Pair.
func (p Pair[A, B]) Size() int {
	return 2
}

// Get returns the value at the given index.
func (p Pair[A, B]) Get(i int) (any, error) {
	switch i {
	case 0:
		return p.first, nil
	case 1:
		return p.second, nil
	default:
		return nil, errors.New("Index out of bounds")
	}
}

// String returns a string representation of the Pair.
func (p Pair[A, B]) String() string {
	return fmt.Sprintf("Pair(%v, %v)", p.first, p.second)
}

// NewPair returns a new Pair instance with the given values.
func NewPair[A, B any](first A, second B) Pair[A, B] {
	return Pair[A, B]{first, second}
}

// Triple is a type that holds three values of generic types.
type Triple[A, B, C any] struct {
	first  A
	second B
	third  C
}

// First returns the first value of the Triple.
func (t Triple[A, B, C]) First() A {
	return t.first
}

// Second returns the second value of the Triple.
func (t Triple[A, B, C]) Second() B {
	return t.second
}

// Third returns the third value of the Triple.
func (t Triple[A, B, C]) Third() C {
	return t.third
}

// Size returns the size of the Triple.
func (t Triple[A, B, C]) Size() int {
	return 3
}

// Get returns the value at the given index.
func (t Triple[A, B, C]) Get(i int) (any, error) {
	switch i {
	case 0:
		return t.first, nil
	case 1:
		return t.second, nil
	case 2:
		return t.third, nil
	default:
		return nil, errors.New("Index out of bounds")
	}
}

// String returns a string representation of the Pair.
func (t Triple[A, B, C]) String() string {
	return fmt.Sprintf("Tripe(%v, %v, %v)", t.first, t.second, t.third)
}

// Third returns the third value of the Triple.
func NewTriple[A, B, C any](first A, second B, third C) Triple[A, B, C] {
	return Triple[A, B, C]{first, second, third}
}
