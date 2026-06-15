package jmespath

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type token struct {
	tokType  tokType
	value    string
	position int
	length   int
}

type tokType int

const eof = -1

// Lexer contains information about the expression being tokenized.
type Lexer struct {
	expression string // The expression provided by the user.
	current    int    // The current position in the string.
	lastWidth  int    // The width of the current rune.  This is used to unread a rune.
	_          []int  // A token buffer.  Not currently used.
}

// SyntaxError is the main error type return when the user provides a bad
// expression.
type SyntaxError struct {
	msg        string // Error message displayed to user
	Expression string // Expression that generated a SyntaxError
	Offset     int    // The location in the string where the error occurred
}

func (e SyntaxError) Error() string {
	// In the error message, we show the caret at the point in the expression
	// where the error occurred.
	hlength := e.Offset
	if hlength < 0 {
		hlength = 0
	}
	return "SyntaxError: " + e.msg + "\n" + e.Expression + "\n" + strings.Repeat(" ", hlength) + "^"
}

func (e SyntaxError) Expression() string {
	return e.Expression
}

func (e SyntaxError) Offset() int {
	return e.Offset
}

type charTypes [utf8.MaxRune + 1]bool

func charTypesFromSlice(s []rune) charTypes {
	var t charTypes
	for _, c := range s {
		t[c] = true
	}
	return t
}

var whiteSpace = charTypesFromSlice([]rune{'\t', '\n', '\r', ' '})
var simpleExpressions = charTypesFromSlice([]rune{'<', '>', '!', '=', '.', '*', ']', ',', ':', '@', '&', '(', ')', '{', '}', '[', '|', '$'})

// NewLexer creates a new JMESPath lexer.
func NewLexer() *Lexer {
	lexer := Lexer{}
	return &lexer
}

// tokenize takes an expression and returns a list of tokens or an error.
// An error is returned if the expression contains invalid syntax.
func (lexer *Lexer) tokenize(expression string) ([]token, error) {
	var tokens []token
	lexer.new(expression)
	for {
		t, err := lexer.nextToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
		if t.tokType == tEOF {
			return tokens, nil
		}
	}
}

// new initializes the lexer with a new expression.
func (lexer *Lexer) new(expression string) {
	lexer.expression = expression
	lexer.current = 0
}

func (lexer *Lexer) nextToken() (token, error) {
	var t token
	for lexer.current < len(lexer.expression) {
		r := lexer.expression[lexer.current]
		if r == ' ' || r == '\n' || r == '\t' || r == '\r' {
			lexer.current++
			continue
		}
		if r >= '0' && r <= '9' {
			t, err := lexer.consumeNumber()
			if err != nil {
				return token{}, err
			}
			return t, nil
		}
		if isAlpha(r) {
			return lexer.consumeUnquotedIdentifier(), nil
		}
		switch r {
		case '.':
			if lexer.current < len(lexer.expression)-1 {
				next := lexer.expression[lexer.current+1]
				if next == '.' {
					t = token{tokType: tFlatten, value: "..", position: lexer.current, length: 2}
					lexer.current += 2
					return t, nil
				}
			}
			t = token{tokType: tDot, value: ".", position: lexer.current, length: 1}
			lexer.current++
		case '*':
			t = token{tokType: tStar, value: "*", position: lexer.current, length: 1}
			lexer.current++
		case ']':
			t = token{tokType: tRbracket, value: "]", position: lexer.current, length: 1}
			lexer.current++
		case ',':
			t = token{tokType: tComma, value: ",", position: lexer.current, length: 1}
			lexer.current++
		case ':':
			t = token{tokType: tColon, value: ":", position: lexer.current, length: 1}
			lexer.current++
		case '@':
			t = token{tokType: tCurrent, value: "@", position: lexer.current, length: 1}
			lexer.current++
		case '&':
			t = token{tokType: tExpref, value: "&", position: lexer.current, length: 1}
			lexer.current++
		case '(':
			t = token{tokType: tLparen, value: "(", position: lexer.current, length: 1}
			lexer.current++
		case ')':
			t = token{tokType: tRparen, value: ")", position: lexer.current, length: 1}
			lexer.current++
		case '{':
			t = token{tokType: tLbrace, value: "{", position: lexer.current, length: 1}
			lexer.current++
		case '}':
			t = token{tokType: tRbrace, value: "}", position: lexer.current, length: 1}
			lexer.current++
		case '[':
			if lexer.current < len(lexer.expression)-1 {
				next := lexer.expression[lexer.current+1]
				if next == ']' {
					t = token{tokType: tFlatten, value: "[]", position: lexer.current, length: 2}
					lexer.current += 2
					return t, nil
				} else if next == '?' {
					t = token{tokType: tFilter, value: "[?", position: lexer.current, length: 2}
					lexer.current += 2
					return t, nil
				}
			}
			t = token{tokType: tLbracket, value: "[", position: lexer.current, length: 1}
			lexer.current++
		case '|':
			t = token{tokType: tPipe, value: "|", position: lexer.current, length: 1}
			lexer.current++
			if lexer.current < len(lexer.expression) && lexer.expression[lexer.current] == '|' {
				t = token{tokType: tOr, value: "||", position: lexer.current - 1, length: 2}
				lexer.current++
			}
		case '!':
			if lexer.current < len(lexer.expression)-1 {
				next := lexer.expression[lexer.current+1]
				if next == '=' {
					t = token{tokType: tNE, value: "!=", position: lexer.current, length: 2}
					lexer.current += 2
					return t, nil
				}
			}
			return token{}, SyntaxError{
				msg:        "Unknown character in expression",
				Expression: lexer.expression,
				Offset:     lexer.current,
			}
		case '<':
			t = token{tokType: tLT, value: "<", position: lexer.current, length: 1}
			lexer.current++
			if lexer.current < len(lexer.expression) && lexer.expression[lexer.current] == '=' {
				t = token{tokType: tLTE, value: "<=", position: lexer.current - 1, length: 2}
				lexer.current++
			}
		case '>':
			t = token{tokType: tGT, value: ">", position: lexer.current, length: 1}
			lexer.current++
			if lexer.current < len(lexer.expression) && lexer.expression[lexer.current] == '=' {
				t = token{tokType: tGTE, value: ">=", position: lexer.current - 1, length: 2}
				lexer.current++
			}
		case '=':
			if lexer.current < len(lexer.expression)-1 {
				next := lexer.expression[lexer.current+1]
				if next == '=' {
					t = token{tokType: tEQ, value: "==", position: lexer.current, length: 2}
					lexer.current += 2
					return t, nil
				}
			}
			return token{}, SyntaxError{
				msg:        "Unknown character in expression: unescaped equals sign",
				Expression: lexer.expression,
				Offset:     lexer.current,
			}
		case '$':
			return token{}, SyntaxError{
				msg:        "Unknown character in expression: dollar sign",
				Expression: lexer.expression,
				Offset:     lexer.current,
			}
		case '"':
			t, err := lexer.consumeQuotedIdentifier()
			if err != nil {
				return t, err
			}
			return t, nil
		case '`':
			t, err := lexer.consumeLiteral()
			if err != nil {
				return t, err
			}
			return t, nil
		case '-':
			t, err := lexer.consumeNumber()
			if err != nil {
				return token{}, err
			}
			return t, nil
		default:
			return token{}, SyntaxError{
				msg:        fmt.Sprintf("Unknown character in expression: %s", string(r)),
				Expression: lexer.expression,
				Offset:     lexer.current,
			}
		}
		return t, nil
	}
	return token{tokType: tEOF, value: "", position: 0, length: 0}, nil
}

func (lexer *Lexer) consumeUnquotedIdentifier() token {
	start := lexer.current
	lexer.current++
	for lexer.current < len(lexer.expression) {
		r := lexer.expression[lexer.current]
		if isAlpha(r) || r >= '0' && r <= '9' {
			lexer.current++
		} else {
			break
		}
	}
	value := lexer.expression[start:lexer.current]
	return token{
		tokType:  tUnquotedIdentifier,
		value:    value,
		position: start,
		length:   lexer.current - start,
	}
}

func (lexer *Lexer) consumeQuotedIdentifier() (token, error) {
	start := lexer.current
	value, err := lexer.consumeStringValue('"')
	if err != nil {
		return token{}, err
	}
	var decoded string
	asJSON := []byte(fmt.Sprintf(`"%s"`, value))
	if err := json.Unmarshal(asJSON, &decoded); err != nil {
		return token{}, SyntaxError{
			msg:        fmt.Sprintf("Invalid escape sequence: %s", err.Error()),
			Expression: lexer.expression,
			Offset:     start,
		}
	}
	return token{
		tokType:  tQuotedIdentifier,
		value:    decoded,
		position: start,
		length:   lexer.current - start,
	}, nil
}

func (lexer *Lexer) consumeLiteral() (token, error) {
	start := lexer.current
	value, err := lexer.consumeStringValue('`')
	if err != nil {
		return token{}, err
	}
	var deserializedValue interface{}
	// To match the JMESPath spec. Literals must be surrounded with backtick characters
	// so when we read the literal, we remove the surrounding backticks.  We then try
	// to JSON parse the literal.  If it fails, it is interpreted as a literal string.
	if err := json.Unmarshal([]byte(value), &deserializedValue); err != nil {
		// Wrap in a double quote string.
		if jsonStr, err := json.Marshal(value); err == nil {
			if err := json.Unmarshal(jsonStr, &deserializedValue); err != nil {
				return token{}, SyntaxError{
					msg:        fmt.Sprintf("Invalid literal: %s", err.Error()),
					Expression: lexer.expression,
					Offset:     start,
				}
			}
		}
	}
	return token{
		tokType:  tJSONLiteral,
		value:    deserializedValue,
		position: start,
		length:   lexer.current - start,
	}, nil
}

func (lexer *Lexer) consumeNumber() (token, error) {
	start := lexer.current
	lexer.current++
	for lexer.current < len(lexer.expression) {
		r := lexer.expression[lexer.current]
		if r >= '0' && r <= '9' {
			lexer.current++
		} else {
			break
		}
	}
	value := lexer.expression[start:lexer.current]
	parsedInt, err := strconv.Atoi(value)
	if err != nil {
		return token{}, SyntaxError{
			msg:        fmt.Sprintf("Invalid number '%s'", value),
			Expression: lexer.expression,
			Offset:     start,
		}
	}
	return token{
		tokType:  tNumber,
		value:    parsedInt,
		position: start,
		length:   lexer.current - start,
	}, nil
}

func (lexer *Lexer) consumeStringValue(delimiter byte) (string, error) {
	start := lexer.current
	lexer.current++
	buffer := bytes.NewBuffer([]byte{})
	for lexer.current < len(lexer.expression) {
		r := lexer.expression[lexer.current]
		if r == '\\' && delimiter == '"' {
			// Escape sequence
			lexer.current++
			if lexer.current == len(lexer.expression) {
				return "", SyntaxError{
					msg:        "Invalid escape sequence (end of expression)",
					Expression: lexer.expression,
					Offset:     start,
				}
			}
			next := lexer.expression[lexer.current]
			if next == '\\' {
				buffer.WriteByte('\\')
				buffer.WriteByte('\\')
			} else {
				buffer.WriteByte('\\')
				buffer.WriteByte(next)
			}
			lexer.current++
		} else if r == delimiter {
			lexer.current++
			return buffer.String(), nil
		} else {
			buffer.WriteByte(r)
			lexer.current++
		}
	}
	return "", SyntaxError{
		msg:        "Unclosed delimiter: " + string(delimiter),
		Expression: lexer.expression,
		Offset:     start,
	}
}

func isAlpha(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || b == '_'
}