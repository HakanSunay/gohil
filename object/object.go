package object

import (
	"fmt"
	"strings"

	"github.com/HakanSunay/gohil/syntaxtree"
)

type Type string

const (
	IntegerObject     Type = "Integer"
	BooleanObject     Type = "Boolean"
	NullObject        Type = "Null"
	ReturnValueObject Type = "ReturnValue"
	ErrorObject       Type = "Error"
	FunctionObject    Type = "Function"
	StringObject      Type = "String"
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

type Function struct {
	Parameters []*syntaxtree.Identifier
	Body       *syntaxtree.BlockStmt
	Env        *Environment
}

func (f *Function) Type() Type {
	return FunctionObject
}

func (f *Function) Inspect() string {
	var builder strings.Builder

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	builder.WriteString("fn")
	builder.WriteString("(")
	builder.WriteString(strings.Join(params, ", "))
	builder.WriteString(") {\n")
	builder.WriteString(f.Body.String())
	builder.WriteString("\n}")

	return builder.String()
}

type String struct {
	Value string
}

func (s *String) Type() Type {
	return StringObject
}

func (s *String) Inspect() string {
	return s.Value
}
