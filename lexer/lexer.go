package lexer

import (
	"unicode"

	"gohil/token"
)

// Lexer is used for lexing the input, which results in producing tokens
// that will later be consumed by the parser (console / interactive mode / repl / shell)
type Lexer struct {
	input string

	// supports only ASCII characters,
	// for Unicode & UTF-8 we need to use rune,
	// but that requires more complex logic for reading multi-byte chars
	currentChar byte

	currentIndex int
	nextIndex    int
}

// New initializes a new lexer type
func New(input string) *Lexer {
	l := &Lexer{input: input}

	// this will initialize the fields of the lexer
	l.nextChar()

	return l
}

// nextChar tries to read the next character for the input field,
// if that is possible, it will update the rest of the fields accordingly
func (l *Lexer) nextChar() {
	// check if next byte is readable
	if l.nextIndex >= len(l.input) {
		// ASCII for NUL
		l.currentChar = 0
		return
	}

	// read next char and update field values
	l.currentChar = l.input[l.nextIndex]
	l.currentIndex = l.nextIndex
	l.nextIndex++
}

// NextToken goes through the input string and extracts the tokens from it
func (l *Lexer) NextToken() token.Token {
	currentToken := token.Token{}

	// major workaround, we are skipping whitespaces
	// for languages like Python, they are necessary for scope definitions
	l.eatWhitespace()

	if unicode.IsDigit(rune(l.currentChar)) {
		currentToken.Type = token.Int
		currentToken.Literal = l.readNumber()

		return currentToken
	}

	if isLetter(l.currentChar) {
		currentToken.Literal = l.readIdentifier()

		// we must wee what kind of an identifier we have read:
		// keyword or user defined
		currentToken.Type = token.ParseIdentifier(currentToken.Literal)

		// no need to move the next byte, since readIdentifier already did that,
		// therefore we can directly return here
		return currentToken
	}

	switch l.currentChar {
	case 0:
		currentToken.Set(token.EOF, l.currentChar)
	case '=':
		if l.peekNextChar() == '=' {
			ch := l.currentChar
			l.nextChar()
			currentToken.Type = token.Equal
			currentToken.Literal = string(ch) + string(l.currentChar)
		} else {
			currentToken.Set(token.Assign, l.currentChar)
		}
	case '+':
		currentToken.Set(token.Plus, l.currentChar)
	case '-':
		currentToken.Set(token.Minus, l.currentChar)
	case '!':
		// not equal
		if l.peekNextChar() == '=' {
			ch := l.currentChar
			l.nextChar()
			currentToken.Type = token.NotEqual
			currentToken.Literal = string(ch) + string(l.currentChar)
		} else {
			currentToken.Set(token.ExclamationMark, l.currentChar)
		}
	case '/':
		currentToken.Set(token.Slash, l.currentChar)
	case '*':
		currentToken.Set(token.Asterisk, l.currentChar)
	case '<':
		currentToken.Set(token.LessThan, l.currentChar)
	case '>':
		currentToken.Set(token.GreaterThan, l.currentChar)
	case ';':
		currentToken.Set(token.SemiColon, l.currentChar)
	case '(':
		currentToken.Set(token.LeftParenthesis, l.currentChar)
	case ')':
		currentToken.Set(token.RightParenthesis, l.currentChar)
	case ',':
		currentToken.Set(token.Comma, l.currentChar)
	case '{':
		currentToken.Set(token.LeftBrace, l.currentChar)
	case '}':
		currentToken.Set(token.RightBrace, l.currentChar)
	case '[':
		currentToken.Set(token.LeftBracket, l.currentChar)
	case ']':
		currentToken.Set(token.RightBracket, l.currentChar)
	case ':':
		currentToken.Set(token.Colon, l.currentChar)
	case '"':
		currentToken.Type = token.String
		currentToken.Literal = l.readString()
	default:
		currentToken.Set(token.Illegal, l.currentChar)
	}

	l.nextChar()
	return currentToken
}

// eatWhitespace is found in a lot of parsers.
// Mostly known as (eat/consume/skip/ignore)Whitespace
// We have decided to rely on the unicode library to choose the actual characters for us.
// Ideally, we might need to define a set of whitespace characters that we might want to skip.
// Some interpreters have tokens for newline characters as well, but we will skip that.
func (l *Lexer) eatWhitespace() {
	// keep on reading till non-whitespace character is hit
	for unicode.IsSpace(rune(l.currentChar)) {
		l.nextChar()
	}
}

// readNumber currently supports only integer values.
// This can be extended to support floats, hexadecimal and octal notation,
// or even complex numbers like GoLang does (1 + 4i).
func (l *Lexer) readNumber() string {
	startIndex := l.currentIndex

	// We can even override this IsDigit method to support Roman numerals,
	for unicode.IsDigit(rune(l.currentChar)) {
		l.nextChar()
	}

	return l.input[startIndex:l.currentIndex]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

// readIdentifier reads the identifier in the input field
// example: "foobar"
// will loop through each character of foobar till it reconstructs the whole word
func (l *Lexer) readIdentifier() string {
	startIndex := l.currentIndex

	// if we want to support identifiers such as
	// foo_bar!_!@_!var7
	// we need to override isLetter to accept those as letters.
	//
	// In GoLang identifiers must abide by the following rule:
	// identifier = letter { letter | unicode_digit } .
	// therefore always starting with a letter,
	// whereas letter abides by:
	// letter = unicode_letter | "_" .
	for isLetter(l.currentChar) {
		l.nextChar()
	}

	return l.input[startIndex:l.currentIndex]
}

// peekNextChar takes a look at the next char in the input string
func (l *Lexer) peekNextChar() byte {
	if l.nextIndex >= len(l.input) {
		return 0
	}

	return l.input[l.nextIndex]
}

// readString reads the whole string starting and ending with (")
// "....."
func (l *Lexer) readString() string {
	startIndex := l.currentIndex + 1
	for {
		l.nextChar()
		if l.currentChar == '"' {
			break
		}
	}

	return l.input[startIndex:l.currentIndex]
}
