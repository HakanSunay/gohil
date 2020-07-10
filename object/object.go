package object

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/HakanSunay/gohil/syntaxtree"
)

type (
	BuiltinFunction func(args ...Object) Object
	Type            string
)

const (
	IntegerObject     Type = "Integer"
	BooleanObject     Type = "Boolean"
	NullObject        Type = "Null"
	ReturnValueObject Type = "ReturnValue"
	ErrorObject       Type = "Error"
	FunctionObject    Type = "Function"
	StringObject      Type = "String"
	BuiltinObject     Type = "Builtin"
	ArrayObject       Type = "Array"
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

// HashKey is used when we are using Boolean objects as keys for Hash Objects
func (i *Integer) HashKey() HashKey {
	// No need to guard for 0 and 1, since they are also used in Boolean.
	// because when we are looking for the exact key hash, we will also use Type,
	// therefore (Type & HashKey) are always equal
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
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

// HashKey is used when we are using Boolean objects as keys for Hash Objects
func (b *Boolean) HashKey() HashKey {
	var hashValue uint64

	if b.Value {
		hashValue = 1
	}

	return HashKey{Type: b.Type(), Value: hashValue}
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

// HashKey is used when we are using String objects as keys for Hash Objects
func (s *String) HashKey() HashKey {
	hash64 := fnv.New64()
	_, err := hash64.Write([]byte(s.Value))
	if err != nil {
		return HashKey{}
	}

	return HashKey{Type: s.Type(), Value: hash64.Sum64()}
}

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() Type {
	return BuiltinObject
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}

type Array struct {
	Elements []Object
}

func (ao *Array) Type() Type {
	return ArrayObject
}

func (ao *Array) Inspect() string {
	var builder strings.Builder

	var elements []string
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	builder.WriteString("[")
	builder.WriteString(strings.Join(elements, ", "))
	builder.WriteString("]")

	return builder.String()
}

type HashKey struct {
	Type  Type
	Value uint64
}
