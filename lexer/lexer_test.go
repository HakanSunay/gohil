package lexer

import (
	"testing"

	"github.com/HakanSunay/gohil/token"
)

func TestLexerNextToken(t *testing.T) {
	type args struct {
		expectedTokenType    token.Type
		expectedTokenLiteral string
	}

	tests := []struct {
		inputString string
		tokenValues []args
	}{
		{
			inputString: `let systemsProgrammingMark = 6;`,
			tokenValues: []args{
				{expectedTokenType: token.Let, expectedTokenLiteral: "let"},
				{expectedTokenType: token.Identifier, expectedTokenLiteral: "systemsProgrammingMark"},
				{expectedTokenType: token.Assign, expectedTokenLiteral: "="},
				{expectedTokenType: token.Int, expectedTokenLiteral: "6"},
				{expectedTokenType: token.SemiColon, expectedTokenLiteral: ";"},
			},
		},

		{
			inputString: `let markCount = 1;`,
			tokenValues: []args{
				{expectedTokenType: token.Let, expectedTokenLiteral: "let"},
				{expectedTokenType: token.Identifier, expectedTokenLiteral: "markCount"},
				{expectedTokenType: token.Assign, expectedTokenLiteral: "="},
				{expectedTokenType: token.Int, expectedTokenLiteral: "1"},
				{expectedTokenType: token.SemiColon, expectedTokenLiteral: ";"},
			},
		},

		{
			inputString: `let calculateMark = fn(markSum, markCount) {
								markSum / markCount;
                          };`,
			tokenValues: []args{
				{expectedTokenType: token.Let, expectedTokenLiteral: "let"},
				{expectedTokenType: token.Identifier, expectedTokenLiteral: "calculateMark"},
				{expectedTokenType: token.Assign, expectedTokenLiteral: "="},
				{expectedTokenType: token.Function, expectedTokenLiteral: "fn"},
				{expectedTokenType: token.LeftParenthesis, expectedTokenLiteral: "("},
				{expectedTokenType: token.Identifier, expectedTokenLiteral: "markSum"},
				{expectedTokenType: token.Comma, expectedTokenLiteral: ","},
				{expectedTokenType: token.Identifier, expectedTokenLiteral: "markCount"},
				{expectedTokenType: token.RightParenthesis, expectedTokenLiteral: ")"},
				{expectedTokenType: token.LeftBrace, expectedTokenLiteral: "{"},
				{expectedTokenType: token.Identifier, expectedTokenLiteral: "markSum"},
				{expectedTokenType: token.Slash, expectedTokenLiteral: "/"},
				{expectedTokenType: token.Identifier, expectedTokenLiteral: "markCount"},
				{expectedTokenType: token.SemiColon, expectedTokenLiteral: ";"},
				{expectedTokenType: token.RightBrace, expectedTokenLiteral: "}"},
				{expectedTokenType: token.SemiColon, expectedTokenLiteral: ";"},
			},
		},

		{
			inputString: `let result = calculateMark(systemsProgrammingMark, markCount);`,
			tokenValues: []args{
				{expectedTokenType: token.Let, expectedTokenLiteral: "let"},
				{expectedTokenType: token.Identifier, expectedTokenLiteral: "result"},
				{expectedTokenType: token.Assign, expectedTokenLiteral: "="},
				{expectedTokenType: token.Identifier, expectedTokenLiteral: "calculateMark"},
				{expectedTokenType: token.LeftParenthesis, expectedTokenLiteral: "("},
				{expectedTokenType: token.Identifier, expectedTokenLiteral: "systemsProgrammingMark"},
				{expectedTokenType: token.Comma, expectedTokenLiteral: ","},
				{expectedTokenType: token.Identifier, expectedTokenLiteral: "markCount"},
				{expectedTokenType: token.RightParenthesis, expectedTokenLiteral: ")"},
				{expectedTokenType: token.SemiColon, expectedTokenLiteral: ";"},
			},
		},

		{
			inputString: `if (result < 6) {
							  return false;
	  					  } else {
	                          return true;
                          }`,
			tokenValues: []args{
				{expectedTokenType: token.If, expectedTokenLiteral: "if"},
				{expectedTokenType: token.LeftParenthesis, expectedTokenLiteral: "("},
				{expectedTokenType: token.Identifier, expectedTokenLiteral: "result"},
				{expectedTokenType: token.LessThan, expectedTokenLiteral: "<"},
				{expectedTokenType: token.Int, expectedTokenLiteral: "6"},
				{expectedTokenType: token.RightParenthesis, expectedTokenLiteral: ")"},
				{expectedTokenType: token.LeftBrace, expectedTokenLiteral: "{"},
				{expectedTokenType: token.Return, expectedTokenLiteral: "return"},
				{expectedTokenType: token.False, expectedTokenLiteral: "false"},
				{expectedTokenType: token.SemiColon, expectedTokenLiteral: ";"},
				{expectedTokenType: token.RightBrace, expectedTokenLiteral: "}"},
				{expectedTokenType: token.Else, expectedTokenLiteral: "else"},
				{expectedTokenType: token.LeftBrace, expectedTokenLiteral: "{"},
				{expectedTokenType: token.Return, expectedTokenLiteral: "return"},
				{expectedTokenType: token.True, expectedTokenLiteral: "true"},
				{expectedTokenType: token.SemiColon, expectedTokenLiteral: ";"},
				{expectedTokenType: token.RightBrace, expectedTokenLiteral: "}"},
			},
		},

		{
			inputString: `"excellent"`,
			tokenValues: []args{
				{expectedTokenType: token.String, expectedTokenLiteral: "excellent"},
			},
		},

		{
			inputString: `[8, 1, 4, 0, 6]`,
			tokenValues: []args{
				{expectedTokenType: token.LeftBracket, expectedTokenLiteral: "["},
				{expectedTokenType: token.Int, expectedTokenLiteral: "8"},
				{expectedTokenType: token.Comma, expectedTokenLiteral: ","},
				{expectedTokenType: token.Int, expectedTokenLiteral: "1"},
				{expectedTokenType: token.Comma, expectedTokenLiteral: ","},
				{expectedTokenType: token.Int, expectedTokenLiteral: "4"},
				{expectedTokenType: token.Comma, expectedTokenLiteral: ","},
				{expectedTokenType: token.Int, expectedTokenLiteral: "0"},
				{expectedTokenType: token.Comma, expectedTokenLiteral: ","},
				{expectedTokenType: token.Int, expectedTokenLiteral: "6"},
				{expectedTokenType: token.RightBracket, expectedTokenLiteral: "]"},
			},
		},

		{
			inputString: `{"81406": "excellent"}`,
			tokenValues: []args{
				{expectedTokenType: token.LeftBrace, expectedTokenLiteral: "{"},
				{expectedTokenType: token.String, expectedTokenLiteral: "81406"},
				{expectedTokenType: token.Colon, expectedTokenLiteral: ":"},
				{expectedTokenType: token.String, expectedTokenLiteral: "excellent"},
				{expectedTokenType: token.RightBrace, expectedTokenLiteral: "}"},
			},
		},

		{
			inputString: `result == 6;`,
			tokenValues: []args{
				{expectedTokenType: token.Identifier, expectedTokenLiteral: "result"},
				{expectedTokenType: token.Equal, expectedTokenLiteral: "=="},
				{expectedTokenType: token.Int, expectedTokenLiteral: "6"},
				{expectedTokenType: token.SemiColon, expectedTokenLiteral: ";"},
			},
		},

		{
			inputString: `result != 2;`,
			tokenValues: []args{
				{expectedTokenType: token.Identifier, expectedTokenLiteral: "result"},
				{expectedTokenType: token.NotEqual, expectedTokenLiteral: "!="},
				{expectedTokenType: token.Int, expectedTokenLiteral: "2"},
				{expectedTokenType: token.SemiColon, expectedTokenLiteral: ";"},
			},
		},

		{
			inputString: ``,
			tokenValues: []args{
				{expectedTokenType: token.EOF, expectedTokenLiteral: string(byte(0))},
			},
		},

		{
			inputString: `?+?-!*>`,
			tokenValues: []args{
				{expectedTokenType: token.Illegal, expectedTokenLiteral: "?"},
				{expectedTokenType: token.Plus, expectedTokenLiteral: "+"},
				{expectedTokenType: token.Illegal, expectedTokenLiteral: "?"},
				{expectedTokenType: token.Minus, expectedTokenLiteral: "-"},
				{expectedTokenType: token.ExclamationMark, expectedTokenLiteral: "!"},
				{expectedTokenType: token.Asterisk, expectedTokenLiteral: "*"},
				{expectedTokenType: token.GreaterThan, expectedTokenLiteral: ">"},
			},
		},

		{
			inputString: `=`,
			tokenValues: []args{
				{expectedTokenType: token.Assign, expectedTokenLiteral: "="},
			},
		},
	}

	for _, tt := range tests {
		l := New(tt.inputString)
		for _, tv := range tt.tokenValues {
			tok := l.NextToken()
			if tok.Type != tv.expectedTokenType {
				t.Errorf("wrong token type, expected: (%s), but got (%s) in [%s]", tv.expectedTokenType, tok.Type, tt.inputString)
			}
			if tok.Literal != tv.expectedTokenLiteral {
				t.Errorf("wrong token literal, expected: (%s), but got: (%s) in [%s]", tv.expectedTokenLiteral, tok.Literal, tt.inputString)
			}
		}
	}
}
