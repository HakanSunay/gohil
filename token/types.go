package token

// Token
type Token struct {
	Type Type
	Literal string
}

// Type identifies the token type and is just an alias to string type.
// String type is easier to work with when debugging,
// but in the ideal case scenario int or byte should be used for performance.
type Type string

const (
	// Illegal is used to identify illegal (unknown) token types
	Illegal = Type("Illegal")

	// EOF is used to identify end of file, the parser can stop after this
	EOF     = Type("EOF")

	// Identifier is used for use inputs such as x, y, z ...
	Identifier = Type("Identifier")

	Int        = Type("Int")

	String     = Type("String")

	// Operators
	Assign      = Type("=")
	Plus        = Type("+")
	Minus       = Type("-")
	Bang        = Type("!")
	Asterisk    = Type("*")
	Slash       = Type("/")
	LessThan    = Type("<")
	GreaterThan = Type(">")
	Equal       = Type("==")
	NotEqual    = Type("!=")

	// Delimiters
	Comma            = Type(",")
	SemiColon        = Type(";")
	LeftParenthesis  = Type("(")
	RightParenthesis = Type(")")
	LeftBrace        = Type("{")
	RightBrace       = Type("}")
	LeftBracket      = Type("[")
	RightBracket     = Type("]")
	Colon            = Type(":")

	// Keywords
	Function = Type("Function")
	Let      = Type("Let")
	True     = Type("True")
	False    = Type("False")
	If       = Type("If")
	Else     = Type("Else")
	Return   = Type("Return")
)

// keywords that are supported by gohil
var keywords = map[string]Type{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

// Parse is used to parse a string to a token type.
// For user inputs that are not in gohil's keywords,
// the result should be of Identifier value
func Parse(input string) Type {
	if val, ok := keywords[input]; ok {
		return val
	}

	return Identifier
}
