package eval

import "github.com/HakanSunay/gohil/object"

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
}
