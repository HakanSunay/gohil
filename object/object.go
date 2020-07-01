package object

import (
	"fmt"
)

type Type string

const (
	IntegerObject     Type = "Integer"
	BooleanObject     Type = "Boolean"
	NullObject        Type = "Null"
	ReturnValueObject Type = "ReturnValue"
	ErrorObject       Type = "Error"
)

type Object interface {
	Type() Type
	Inspect() string // String repr
}

type Integer struct {
	Value int
}

func (i *Integer) Type() Type {
	return IntegerObject
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type {
	return BooleanObject
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Null references were introduced to the ALGOL W language in 1965
// Tony Hoare was the person behind this and later on he called this
// a “billion-dollar mistake”
type Null struct{}

func (n *Null) Type() Type {
	return NullObject
}

func (n *Null) Inspect() string {
	return "null"
}

// ReturnValue is used as a wrapper around the to-be returned Value Object.
// This is strictly done to skip doing ugly go to statements.
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() Type {
	return ReturnValueObject
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() Type {
	return ErrorObject
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}
