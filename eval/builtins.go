package eval

import (
	"fmt"

	"github.com/HakanSunay/gohil/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			// generally len works with 1 argument only :)
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			// if the fist argument is a string, just find its len using the GoLang len keyword
			case *object.String:
				return &object.Integer{Value: len(arg.Value)}
			// we can always add a new case here for custom behaviour for certain object.Object :)
			// make it work for int as well, but what is the LEN of an int? (no one knows, yet :) )
			case *object.Array:
				return &object.Integer{Value: len(arg.Elements)}
			default:
				return newError("argument of `len` not supported, got %s", args[0].Type())
			}
		},
	},
	// Calling this head to remind myself of the painful logical programming days
	"head": {
		Fn: func(args ...object.Object) object.Object {
			// generally head works with 1 argument only :)
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ArrayObject {
				return newError("argument of head must be of type Array, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return Null
		},
	},
	// Prolog analogy
	"tail": {
		Fn: func(args ...object.Object) object.Object {
			// generally head works with 1 argument only :)
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ArrayObject {
				return newError("argument of head must be of type Array, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if length := len(arr.Elements); length > 0 {
				// let's not modify the old object
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[1:length])

				return &object.Array{Elements: newElements}
			}

			return Null
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			// generally head works with 1 argument only :)
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ArrayObject {
				return newError("argument of head must be of type Array, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if length := len(arr.Elements); length > 0 {
				return arr.Elements[length-1]
			}

			return Null
		},
	},
	"append": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			if args[0].Type() != object.ArrayObject {
				return newError("argument of append must be of type Array, got %s", args[0].Type())
			}

			// creating a new object, not modifying the old one
			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)

			// add the new element
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, a := range args {
				// This is given the fact that STDOUT is used,
				// if in any scenario we want to do it any other way around
				// we can also create a object.String with a string builder
				// and after the loop assigning the builder value to a object.String
				// but that could mean that the object can be used somewhere else

				// but print is all about printing to stdout and not producing values!!!
				fmt.Println(a.Inspect())
			}

			return Null
		},
	},
}
